package bootstrap

import (
	"errors"
	"io/fs"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestBuildTransfersResourcesAfterInitialization(t *testing.T) {
	config := newTestConfig(t)
	var calls []string

	if err := config.Resources.Own("test", func() error {
		calls = append(calls, "resource")

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"first", "second"} {
		config.Modules = append(config.Modules, registrationModule(func(boot module.Bootstrap) error {
			calls = append(calls, "register "+name)
			boot.Hooks().Engine().OnInit(func() error {
				calls = append(calls, "init "+name)

				return nil
			})
			boot.Hooks().Engine().OnClose(func() error {
				calls = append(calls, "close "+name)

				return nil
			})

			return nil
		}))
	}

	result, err := Build(config)
	if err != nil {
		t.Fatal(err)
	}

	if result.Compiler == nil || result.DebugCompiler == nil || result.Host == nil || result.Hooks == nil || result.Limiter == nil {
		t.Fatalf("missing initialized dependencies: %+v", result)
	}

	if result.Resources != config.Resources || result.Hooks == config.Hooks {
		t.Fatal("build must transfer the resource manager and snapshot the hooks")
	}

	want := []string{"register first", "register second", "init first", "init second"}
	if !slices.Equal(calls, want) {
		t.Fatalf("construction calls = %v, want %v", calls, want)
	}

	if _, err := result.Host.FileSystem.Stat("."); err != nil {
		t.Fatalf("transferred filesystem is unusable: %v", err)
	}

	if err := result.Hooks.EngineHooks.RunClose(); err != nil {
		t.Fatal(err)
	}

	if err := result.Resources.Close(); err != nil {
		t.Fatal(err)
	}

	if err := result.Resources.Close(); err != nil {
		t.Fatal(err)
	}

	want = append(want, "close second", "close first", "resource")
	if !slices.Equal(calls, want) {
		t.Fatalf("ownership transfer cleanup = %v, want %v", calls, want)
	}

	if _, err := result.Host.FileSystem.Stat("."); err == nil {
		t.Fatal("owner cleanup left the transferred filesystem open")
	}
}

func TestBuildRollbackPreservesErrorsAndCleanupOrder(t *testing.T) {
	for _, failingCleanup := range []bool{false, true} {
		t.Run(map[bool]string{false: "successful cleanup", true: "failing cleanup"}[failingCleanup], func(t *testing.T) {
			config := newTestConfig(t)
			registerErr := errors.New("registration failed")
			hookFailure := &fs.PathError{Op: "flush", Path: "hook", Err: errors.New("failed")}
			var hookErr, firstErr, secondErr error

			if failingCleanup {
				hookErr = hookFailure
				firstErr = errors.New("first failed")
				secondErr = errors.New("second failed")
			}

			var calls []string

			if err := config.Resources.Own("first", func() error {
				calls = append(calls, "resource first")

				return firstErr
			}); err != nil {
				t.Fatal(err)
			}

			config.Hooks.EngineHooks.OnClose(func() error {
				calls = append(calls, "close option")

				return nil
			})
			config.Modules = []module.Module{registrationModule(func(boot module.Bootstrap) error {
				boot.Hooks().Engine().OnClose(func() error {
					calls = append(calls, "close module")

					if _, err := boot.Host().FileSystem().Stat("."); err != nil {
						t.Errorf("filesystem closed before hooks: %v", err)
					}

					return hookErr
				})

				if err := config.Resources.Own("second", func() error {
					calls = append(calls, "resource second")

					return secondErr
				}); err != nil {
					t.Fatal(err)
				}

				return registerErr
			})}

			result, err := Build(config)
			if result != (Result{}) || !errors.Is(err, registerErr) {
				t.Fatalf("construction result = %+v/%v, want zero result and registration error", result, err)
			}

			if failingCleanup {
				counts := make(map[error]int)
				pending := []error{err}
				for len(pending) > 0 {
					current := pending[len(pending)-1]
					pending = pending[:len(pending)-1]
					counts[current]++

					switch wrapped := current.(type) {
					case interface{ Unwrap() []error }:
						pending = append(pending, wrapped.Unwrap()...)
					case interface{ Unwrap() error }:
						pending = append(pending, wrapped.Unwrap())
					}
				}

				for _, cause := range []error{registerErr, hookErr, firstErr, secondErr} {
					if !errors.Is(err, cause) {
						t.Fatalf("lost cleanup cause %v: %v", cause, err)
					}

					if counts[cause] != 1 {
						t.Errorf("rollback cause %v occurs %d times, want once", cause, counts[cause])
					}
				}

				var typed *fs.PathError
				if !errors.As(err, &typed) || typed != hookFailure {
					t.Fatalf("lost typed hook error: %v", err)
				}

				const wantError = "registration failed\nclose engine: close hooks: flush hook: failed\nclose second: second failed\nclose first: first failed"
				if err.Error() != wantError {
					t.Errorf("rollback error = %q, want %q", err.Error(), wantError)
				}
			} else if err != registerErr {
				t.Fatalf("successful cleanup wrapped the registration error: %v", err)
			}

			want := []string{"close module", "close option", "resource second", "resource first"}
			if !slices.Equal(calls, want) {
				t.Fatalf("rollback calls = %v, want %v", calls, want)
			}

			_ = config.Resources.Close()

			if !slices.Equal(calls, want) {
				t.Fatalf("repeated resource close retried rollback cleanup: %v", calls)
			}
		})
	}
}

func TestBuildInitFailureRollsBackHookSnapshot(t *testing.T) {
	config := newTestConfig(t)
	initErr := errors.New("initialization failed")
	var calls []string

	if err := config.Resources.Own("test", func() error {
		calls = append(calls, "resource")

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	config.Modules = []module.Module{registrationModule(func(boot module.Bootstrap) error {
		boot.Hooks().Engine().OnClose(func() error {
			calls = append(calls, "close snapshot")

			return nil
		})
		boot.Hooks().Engine().OnInit(func() error {
			calls = append(calls, "init")
			boot.Hooks().Engine().OnClose(func() error {
				calls = append(calls, "close late")

				return nil
			})

			return initErr
		})
		boot.Hooks().Engine().OnInit(func() error {
			calls = append(calls, "init skipped")

			return nil
		})

		return nil
	})}

	result, err := Build(config)
	if result != (Result{}) || !errors.Is(err, initErr) {
		t.Fatalf("construction result = %+v/%v, want zero result and init error", result, err)
	}

	if err.Error() != "init hooks: initialization failed" {
		t.Fatalf("initialization error changed: %v", err)
	}

	want := []string{"init", "close snapshot", "resource"}
	if !slices.Equal(calls, want) {
		t.Fatalf("initialization rollback = %v, want %v", calls, want)
	}
}

func newTestConfig(t *testing.T) Config {
	t.Helper()

	root := t.TempDir()
	resources := resource.NewManager()
	t.Cleanup(func() { _ = resources.Close() })

	return Config{
		Host: host.Config{
			Library:  runtime.NewLibrary(),
			Params:   runtime.NewParams(),
			Encoding: encoding.NewRegistry(),
			FSRoot:   root,
		},
		Hooks:             host.NewHooks(),
		Resources:         resources,
		OptimizationLevel: compiler.Full,
		MaxActiveSessions: 1,
	}
}
