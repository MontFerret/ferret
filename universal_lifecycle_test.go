package ferret

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

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
