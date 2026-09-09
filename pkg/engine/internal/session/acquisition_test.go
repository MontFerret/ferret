package session

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestAcquisitionRollsBackPartialFailures(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, stage := range []string{"logger", "filesystem", "environment", "environment panic", "vm", "debug execution", "debug session"} {
			if !debug && strings.HasPrefix(stage, "debug") {
				continue
			}

			t.Run(map[bool]string{false: "ordinary/", true: "debug/"}[debug]+stage, func(t *testing.T) {
				f := newSessionFixture(t, 1)
				config := NewConfig()
				config.FSRoot = t.TempDir()
				hooks := 0
				f.hooks.OnClose(func() error {
					hooks++

					return nil
				})
				panicValue := errors.New("environment panic")
				switch stage {
				case "logger":
					config.Logger = []logging.Option{logging.WithLevel(logging.LogLevel(100))}
				case "filesystem":
					config.FSRoot = filepath.Join(t.TempDir(), "missing")
				case "environment":
					config.Environment = []vm.EnvironmentOption{vm.WithParam("", runtime.Int(1))}
				case "environment panic":
					config.Environment = []vm.EnvironmentOption{vm.WithFunctionsRegistrar(func(runtime.FunctionDefs) { panic(panicValue) })}
				case "vm":
					if debug {
						f.program.Registers = 0
					} else if err := f.pool.Close(); err != nil {
						t.Fatal(err)
					}
				case "debug execution":
					f.program.Metadata.DebugPoints = nil
				case "debug session":
					f.program.Source = source.Source{}
				}

				a := acquisition{limiter: f.limiter}
				var resultErr error
				var recovered any
				func() {
					defer func() { recovered = recover() }()
					resultErr = func() (err error) {
						defer a.rollback(&err)
						if err := a.prepare(t.Context(), config, f.host, f.closed); err != nil {
							return err
						}

						if stage == "logger" || stage == "filesystem" || strings.HasPrefix(stage, "environment") {
							t.Fatal("expected preparation failure")
						}

						var child io.Closer
						if debug {
							child, err = a.newDebugSession(config, f.host, f.hooks, f.program)
						} else {
							child, err = a.newExecution(config, f.host, f.hooks, f.pool)
						}

						if err == nil {
							_ = child.Close()
							t.Fatal("expected VM/debugger construction failure")
						}

						return err
					}()
				}()

				if stage == "environment panic" {
					if recovered != panicValue {
						t.Fatalf("panic=%v, want %v", recovered, panicValue)
					}
				} else if recovered != nil || resultErr == nil {
					t.Fatalf("panic=%v error=%v", recovered, resultErr)
				}

				if stage == "debug execution" && !strings.Contains(resultErr.Error(), "program has no debugger metadata") {
					t.Fatalf("failure did not reach debug execution setup: %v", resultErr)
				}

				if stage == "debug session" && !strings.Contains(resultErr.Error(), "debug source is required") {
					t.Fatalf("failure did not reach debugger session setup: %v", resultErr)
				}

				if !debug && stage == "vm" && (!errors.Is(resultErr, runtime.ErrInvalidOperation) || resultErr.Error() != runtime.Error(runtime.ErrInvalidOperation, "plan is closed").Error()) {
					t.Fatalf("closed pool translation changed: %v", resultErr)
				}

				if hooks != 0 || len(f.limiter.ch) != 0 {
					t.Fatalf("partial construction ran hooks or retained capacity: hooks=%d permits=%d", hooks, len(f.limiter.ch))
				}

				if a.filesystem != nil && a.filesystem != f.host.FileSystem {
					if _, err := a.filesystem.Stat("."); err == nil {
						t.Fatal("rollback retained an owned filesystem")
					}
				}

				if _, err := f.host.FileSystem.Stat("."); err != nil {
					t.Fatalf("rollback closed the borrowed host filesystem: %v", err)
				}
			})
		}
	}
}

func TestAcquisitionRollbackPreservesErrorsAndPanicRelease(t *testing.T) {
	for _, mode := range []string{"error", "construction panic", "cleanup panic"} {
		t.Run(mode, func(t *testing.T) {
			f := newSessionFixture(t, 1)
			buildErr, closeErr := errors.New("construction failed"), errors.New("filesystem failed")
			filesystem := newCountingCloseFileSystem(t, t.TempDir(), closeErr)
			var recovered any
			var resultErr error
			func() {
				defer func() { recovered = recover() }()
				resultErr = func() (err error) {
					a := acquisition{limiter: f.limiter}
					defer a.rollback(&err)
					if err := a.prepare(t.Context(), NewConfig(), f.host, f.closed); err != nil {
						return err
					}

					if err := a.resources.Own(resource.FileSystem, func() error {
						err := filesystem.Close()
						if mode == "cleanup panic" {
							panic(closeErr)
						}

						return err
					}); err != nil {
						return err
					}

					if mode == "construction panic" {
						panic(buildErr)
					}

					return buildErr
				}()
			}()

			wantPanic := map[string]error{"construction panic": buildErr, "cleanup panic": closeErr}[mode]
			if recovered != wantPanic {
				t.Fatalf("panic=%v, want %v", recovered, wantPanic)
			}

			if mode == "error" {
				if !errors.Is(resultErr, buildErr) || !errors.Is(resultErr, closeErr) || !strings.Contains(resultErr.Error(), "close filesystem") {
					t.Fatalf("rollback lost construction or cleanup failure: %v", resultErr)
				}

				joined := resultErr.(interface{ Unwrap() []error }).Unwrap()
				if len(joined) != 2 || joined[0] != buildErr || !errors.Is(joined[1], closeErr) {
					t.Fatalf("rollback changed error order: %v", joined)
				}
			}

			if filesystem.closeCalls.Load() != 1 || len(f.limiter.ch) != 0 {
				t.Fatalf("rollback closes=%d permits=%d", filesystem.closeCalls.Load(), len(f.limiter.ch))
			}

			if _, err := filesystem.Stat("."); err == nil {
				t.Fatal("rollback did not close its owned filesystem")
			}

			if _, err := f.host.FileSystem.Stat("."); err != nil {
				t.Fatalf("rollback closed the borrowed filesystem: %v", err)
			}

			f.execution(t, NewConfig())
		})
	}
}

func TestConstructorsReleasePermitAfterEnvironmentPanic(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			f := newSessionFixture(t, 1)
			panicValue := errors.New("environment panic")
			config := NewConfig()
			config.FSRoot = t.TempDir()
			config.Environment = []vm.EnvironmentOption{vm.WithFunctionsRegistrar(func(runtime.FunctionDefs) { panic(panicValue) })}
			hooks := 0
			f.hooks.OnClose(func() error {
				hooks++

				return nil
			})
			var recovered any
			func() {
				defer func() { recovered = recover() }()

				if debug {
					_, _ = NewDebugSession(t.Context(), config, f.host, f.hooks, f.limiter, f.program, f.closed)
				} else {
					_, _ = NewExecution(t.Context(), config, f.host, f.hooks, f.limiter, f.pool, f.closed)
				}
			}()

			if recovered != panicValue || hooks != 0 || len(f.limiter.ch) != 0 {
				t.Fatalf("panic=%v hooks=%d permits=%d", recovered, hooks, len(f.limiter.ch))
			}

			f.execution(t, NewConfig())
		})
	}
}

func TestConstructionTransfersBeforeCancellationCheck(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			f := newSessionFixture(t, 1)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			hookErr, resourceErr := errors.New("close hook"), errors.New("host resource")
			hooks, closes := 0, 0
			f.hooks.OnClose(func() error {
				hooks++

				return hookErr
			})
			config := NewConfig()
			config.FSRoot = t.TempDir()
			a := acquisition{limiter: f.limiter}
			config.Environment = []vm.EnvironmentOption{vm.WithFunctionsRegistrar(func(runtime.FunctionDefs) {
				cancel()
				if err := a.resources.Own("publication-service", func() error {
					closes++

					return resourceErr
				}); err != nil {
					t.Fatal(err)
				}
			})}

			var err error
			defer a.rollback(&err)
			if err := a.prepare(ctx, config, f.host, f.closed); err != nil {
				t.Fatal(err)
			}

			var session io.Closer
			if debug {
				session, err = a.newDebugSession(config, f.host, f.hooks, f.program)
			} else {
				session, err = a.newExecution(config, f.host, f.hooks, f.pool)
			}

			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = session.Close() })
			a.rollback(&err)
			if ctx.Err() != context.Canceled || hooks != 0 || closes != 0 || len(f.limiter.ch) != 1 {
				t.Fatalf("construction rolled back before native publication: hooks=%d closes=%d permits=%d", hooks, closes, len(f.limiter.ch))
			}

			if err := session.Close(); !errors.Is(err, hookErr) || !errors.Is(err, resourceErr) {
				t.Fatalf("completed session lost cleanup causes: %v", err)
			}

			if hooks != 1 || closes != 1 || len(f.limiter.ch) != 0 {
				t.Fatalf("cleanup hooks=%d closes=%d permits=%d", hooks, closes, len(f.limiter.ch))
			}

			if _, err := a.filesystem.Stat("."); err == nil {
				t.Fatal("completed session retained its owned filesystem")
			}
		})
	}
}
