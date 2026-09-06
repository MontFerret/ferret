package ferret

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCompileCancellationAroundHooks(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "normal", true: "debug"}[debug], func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			engine := mustNewEngine(t, WithAfterCompileHook(func(context.Context, error) error { cancel(); return nil }))
			t.Cleanup(func() { _ = engine.Close() })
			compile := engine.Compile

			if debug {
				compile = engine.CompileDebug
			}

			plan, err := compile(ctx, NewAnonymousSource("RETURN 1"))
			if plan != nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("plan=%v err=%v", plan, err)
			}
		})
	}
}

func TestEngineCloseSettlesAdmittedCompilation(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	releaseHook := sync.OnceFunc(func() { close(release) })
	t.Cleanup(releaseHook)
	var cleanup atomic.Int32
	engine := mustNewEngine(t, WithBeforeCompileHook(func(context.Context) error { close(entered); <-release; return nil }), WithEngineCloseHook(func() error { cleanup.Add(1); return nil }))
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)
	result := make(chan error, 1)
	go func() {
		p, err := engine.Compile(ctx, NewAnonymousSource("RETURN 1"))
		if p != nil {
			_ = p.Close()
			err = errors.New("published plan after close")
		}

		result <- err
	}()
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

	stopping := engine.lifecycle.stopped()
	closed := make(chan error, 1)
	go func() { closed <- engine.Close() }()
	select {
	case <-stopping:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

	if cleanup.Load() != 0 {
		t.Fatal("parent cleaned up during admitted compilation")
	}

	releaseHook()
	select {
	case err := <-result:
		if !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

	select {
	case err := <-closed:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

	if cleanup.Load() != 1 {
		t.Fatal("engine cleanup count")
	}
}

func TestParentCloseInterruptsBlockedSessionAdmission(t *testing.T) {
	for _, closeEngine := range []bool{false, true} {
		for _, debug := range []bool{false, true} {
			t.Run(map[bool]string{false: "plan", true: "engine"}[closeEngine]+map[bool]string{false: "/ordinary", true: "/debug"}[debug], func(t *testing.T) {
				engine := mustNewEngine(t, WithMaxActiveSessions(1))

				plan, err := engine.CompileDebug(t.Context(), NewAnonymousSource("RETURN 1"))
				if err != nil {
					t.Fatal(err)
				}

				first, err := plan.NewSession(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = first.Close(); _ = plan.Close(); _ = engine.Close() })
				ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
				t.Cleanup(cancel)
				entered, release := make(chan struct{}), make(chan struct{})
				releaseHook := sync.OnceFunc(func() { close(release) })
				t.Cleanup(releaseHook)
				result := make(chan error, 1)
				option := func(api.SessionOptions) error { close(entered); <-release; return nil }
				go func() {
					if debug {
						session, err := plan.NewDebugSession(ctx, option)
						if session != nil {
							_ = session.Close()
						}

						result <- err
					} else {
						session, err := plan.NewSession(ctx, option)
						if session != nil {
							_ = session.Close()
						}

						result <- err
					}
				}()
				select {
				case <-entered:
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				}

				parent := &plan.lifecycle
				closeParent := plan.Close

				if closeEngine {
					parent = &engine.lifecycle
					closeParent = engine.Close
				}

				stopping := parent.stopped()
				closed := make(chan error, 1)
				go func() { closed <- closeParent() }()
				select {
				case <-stopping:
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				}

				releaseHook()
				select {
				case err := <-result:
					if !errors.Is(err, runtime.ErrInvalidOperation) {
						t.Fatalf("admission=%v", err)
					}
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				}

				select {
				case err := <-closed:
					if err != nil {
						t.Fatal(err)
					}
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				}
			})
		}
	}
}

func TestNativeCloseRetainsResultForConcurrentCallers(t *testing.T) {
	engineErr, planErr := errors.New("engine cleanup"), errors.New("plan cleanup")
	var engines, plans atomic.Int32
	engine := mustNewEngine(t, WithEngineCloseHook(func() error { engines.Add(1); return engineErr }), WithPlanCloseHook(func() error { plans.Add(1); return planErr }))
	plan := mustCompilePlan(t, engine, "RETURN 1")
	for _, operation := range []struct {
		close func() error
		cause error
	}{{plan.Close, planErr}, {engine.Close, engineErr}} {
		var wg sync.WaitGroup
		for range 8 {
			wg.Go(func() {
				if err := operation.close(); !errors.Is(err, operation.cause) {
					t.Errorf("close: %v", err)
				}
			})
		}

		wg.Wait()
	}

	if engines.Load() != 1 || plans.Load() != 1 {
		t.Fatalf("closes=%d/%d", engines.Load(), plans.Load())
	}

	if p, err := engine.Compile(t.Context(), NewAnonymousSource("RETURN 1")); p != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("compile closed engine: %v/%v", p, err)
	}
}

func TestNativeNilAndCanceledContextsRejectAdmission(t *testing.T) {
	engine := mustNewEngine(t)
	t.Cleanup(func() { _ = engine.Close() })

	plan, err := engine.CompileDebug(t.Context(), NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	ordinary, err := plan.NewSession(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = ordinary.Close() })

	debug, err := plan.NewDebugSession(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = debug.Close() })
	operations := []func(context.Context) error{
		func(ctx context.Context) error {
			_, err := engine.Compile(ctx, NewAnonymousSource("RETURN 1"))

			return err
		},
		func(ctx context.Context) error {
			_, err := engine.CompileDebug(ctx, NewAnonymousSource("RETURN 1"))

			return err
		},
		func(ctx context.Context) error { _, err := plan.NewSession(ctx); return err },
		func(ctx context.Context) error { _, err := plan.NewDebugSession(ctx); return err },
		func(ctx context.Context) error { _, err := ordinary.Run(ctx); return err },
		func(ctx context.Context) error { _, err := debug.Start(ctx); return err },
		func(ctx context.Context) error { _, err := debug.Continue(ctx); return err },
		func(ctx context.Context) error { _, err := debug.StepIn(ctx); return err },
		func(ctx context.Context) error { _, err := debug.StepOver(ctx); return err },
		func(ctx context.Context) error { _, err := debug.StepOut(ctx); return err },
		func(ctx context.Context) error { _, err := debug.Evaluate(ctx, "1"); return err },
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for i, operation := range operations {
		if err := operation(nil); !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Errorf("nil operation %d: %v", i, err)
		}

		if err := operation(ctx); !errors.Is(err, context.Canceled) {
			t.Errorf("canceled operation %d: %v", i, err)
		}
	}
}
