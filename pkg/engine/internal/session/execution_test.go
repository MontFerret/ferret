package session

import (
	"errors"
	"io"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestExecutionRunsSequentiallyWithIsolatedEnvironment(t *testing.T) {
	f := newSessionFixture(t, 1)
	config := NewConfig()
	if err := config.SetParam("value", 41); err != nil {
		t.Fatal(err)
	}

	execution := f.execution(t, config)
	for range 2 {
		result, err := execution.Run(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		output, encodeErr, closeErr := execution.MaterializeAndClose(result)
		if encodeErr != nil || closeErr != nil || output == nil || string(output.Content) != "42" {
			t.Fatalf("output=%v encoding=%v cleanup=%v", output, encodeErr, closeErr)
		}

		if result.Root() != runtime.None {
			t.Fatal("materialization retained its VM result")
		}
	}

	if f.host.Params["value"] != runtime.Int(0) {
		t.Fatal("session overrides mutated host parameters")
	}
}

func TestExecutionMaterializeAndCloseSeparatesErrors(t *testing.T) {
	for _, missingCodec := range []bool{false, true} {
		t.Run(map[bool]string{false: "encoded", true: "missing codec"}[missingCodec], func(t *testing.T) {
			f := newSessionFixture(t, 1)
			config := NewConfig()
			if missingCodec {
				config.OutputContentType = "application/missing"
			}

			execution := f.execution(t, config)
			result, err := execution.Run(t.Context())
			if err != nil {
				t.Fatal(err)
			}

			cleanupErr := errors.New("result cleanup")
			filesystem := newCountingCloseFileSystem(t, t.TempDir(), cleanupErr)
			result.AdoptCloser(filesystem)
			output, encodeErr, closeErr := execution.MaterializeAndClose(result)
			if !errors.Is(closeErr, cleanupErr) || filesystem.closeCalls.Load() != 1 {
				t.Fatalf("cleanup=%v closes=%d", closeErr, filesystem.closeCalls.Load())
			}

			if missingCodec {
				if output != nil || !errors.Is(encodeErr, encoding.ErrCodecNotFound) {
					t.Fatalf("output=%v encoding=%v", output, encodeErr)
				}
			} else if output == nil || string(output.Content) != "1" || encodeErr != nil {
				t.Fatalf("output=%v encoding=%v", output, encodeErr)
			}

			if err := execution.Close(); err != nil {
				t.Fatal(err)
			}

			if filesystem.closeCalls.Load() != 1 || result.Root() != runtime.None {
				t.Fatal("session closure repeated or retained result cleanup")
			}
		})
	}
}

func TestSessionCleanupIsConcurrentAndRetainsErrors(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				f := newSessionFixture(t, 1)
				hookErr, closeErr := errors.New("hook failed"), errors.New("service failed")
				var hooks, closes atomic.Int32
				closing, resume := make(chan struct{}), make(chan struct{})
				releaseCleanup := sync.OnceFunc(func() { close(resume) })
				defer releaseCleanup()
				var calls []string
				f.hooks.OnClose(func() error {
					hooks.Add(1)
					calls = append(calls, "hook")

					return hookErr
				})

				config := NewConfig()
				config.FSRoot = t.TempDir()
				a := acquisition{limiter: f.limiter}
				var err error
				defer a.rollback(&err)

				if err := a.prepare(t.Context(), config, f.host, f.closed); err != nil {
					t.Fatal(err)
				}

				filesystem := a.filesystem
				if err := a.resources.Own("service", func() error {
					closes.Add(1)
					calls = append(calls, "resource")
					if hooks.Load() != 1 || len(f.limiter.ch) != 1 {
						t.Error("resource cleanup violated hook/permit ordering")
					}

					close(closing)
					<-resume

					return closeErr
				}); err != nil {
					t.Fatal(err)
				}

				var session io.Closer
				var instance *vm.VM
				if debug {
					session, err = a.newDebugSession(config, f.host, f.hooks, f.program)
				} else {
					var execution *Execution
					execution, err = a.newExecution(config, f.host, f.hooks, f.pool)
					if err == nil {
						instance = execution.vm
					}

					session = execution
				}

				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = session.Close() })
				a.rollback(&err)
				if len(f.limiter.ch) != 1 {
					t.Fatal("successful construction did not transfer permit ownership")
				}

				results := make([]error, 8)
				var wg sync.WaitGroup
				wg.Go(func() { results[0] = session.Close() })
				<-closing
				if len(f.limiter.ch) != 1 {
					t.Fatal("permit released before resources closed")
				}

				if !debug {
					if _, err := f.pool.Acquire(); !errors.Is(err, vm.ErrPoolExhausted) {
						t.Fatalf("VM returned before resource cleanup: %v", err)
					}
				}

				for i := 1; i < len(results); i++ {
					wg.Go(func() { results[i] = session.Close() })
				}

				releaseCleanup()
				wg.Wait()
				for _, err := range results {
					if !errors.Is(err, hookErr) || !errors.Is(err, closeErr) || err != results[0] {
						t.Fatalf("concurrent closure lost cleanup errors: %v", err)
					}
				}

				if hooks.Load() != 1 || closes.Load() != 1 || !slices.Equal(calls, []string{"hook", "resource"}) {
					t.Fatalf("cleanup calls=%v hooks=%d closes=%d", calls, hooks.Load(), closes.Load())
				}

				if _, err := filesystem.Stat("."); err == nil {
					t.Fatal("resource cleanup failure left the filesystem open")
				}

				next := f.execution(t, NewConfig())
				if !debug && next.vm != instance {
					t.Fatal("cleanup failure prevented returning the VM to its pool")
				}

				if err := session.Close(); err != results[0] || len(f.limiter.ch) != 1 {
					t.Fatalf("repeated closure changed error or consumed a sibling permit: %v", err)
				}
			})
		})
	}
}

func TestExecutionCloseReturnsBorrowedVMOnce(t *testing.T) {
	f := newSessionFixture(t, 2)
	execution := f.execution(t, NewConfig())
	instance := execution.vm
	if err := execution.Close(); err != nil {
		t.Fatal(err)
	}

	next := f.execution(t, NewConfig())
	if next.vm != instance {
		t.Fatal("the next execution did not reuse the returned VM")
	}

	if err := execution.Close(); err != nil {
		t.Fatal(err)
	}

	if len(f.limiter.ch) != 1 {
		t.Fatal("repeated close released a sibling permit")
	}

	if _, err := f.pool.Acquire(); !errors.Is(err, vm.ErrPoolExhausted) {
		t.Fatalf("repeated close returned the sibling's VM: %v", err)
	}
}

func TestExecutionCloseAfterPoolClose(t *testing.T) {
	f := newSessionFixture(t, 1)
	execution := f.execution(t, NewConfig())
	instance := execution.vm
	if err := f.pool.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := execution.Run(t.Context())
	if err != nil {
		t.Fatalf("pool closure revoked the borrowed VM: %v", err)
	}

	if _, encodeErr, closeErr := execution.MaterializeAndClose(result); encodeErr != nil || closeErr != nil {
		t.Fatalf("encoding=%v cleanup=%v", encodeErr, closeErr)
	}

	if err := execution.Close(); err != nil || len(f.limiter.ch) != 0 {
		t.Fatalf("close=%v permits=%d", err, len(f.limiter.ch))
	}

	if result, err := instance.Run(t.Context(), vm.NewDefaultEnvironment()); err == nil {
		_ = result.Close()
		t.Fatal("return to closed pool did not close the VM")
	}
}

func TestExecutionCloseJoinsOwnedFileSystemErrorExactlyOnce(t *testing.T) {
	f := newSessionFixture(t, 1)
	hookErr, fsErr := errors.New("hook failed"), errors.New("filesystem failed")
	f.hooks.OnClose(func() error { return hookErr })
	execution := f.execution(t, NewConfig())
	filesystem := newCountingCloseFileSystem(t, t.TempDir(), fsErr)
	if err := execution.resources.Own(resource.FileSystem, filesystem.Close); err != nil {
		t.Fatal(err)
	}

	execution.filesystem = filesystem
	first := execution.Close()
	if !errors.Is(first, hookErr) || !errors.Is(first, fsErr) || first != execution.Close() || filesystem.closeCalls.Load() != 1 {
		t.Fatalf("close=%v calls=%d", first, filesystem.closeCalls.Load())
	}

	if _, err := f.host.FileSystem.Stat("."); err != nil {
		t.Fatalf("closing an owned replacement closed the borrowed host filesystem: %v", err)
	}
}
