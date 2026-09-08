package ferret

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func TestPlanOptionsApplyInOrderAndJoinFailures(t *testing.T) {
	firstErr, lastErr := errors.New("first option"), errors.New("last option")
	var calls []string
	opts, err := newPlanOptions(compiler.Full, false, []PlanOption{
		nil,
		func(*planConfig) error {
			calls = append(calls, "first")

			return firstErr
		},
		WithPlanOptimizationLevel(OptimizationNone),
		func(opts *planConfig) error {
			calls = append(calls, "middle")

			if opts.level != compiler.None {
				t.Errorf("preceding option was not applied: %v", opts.level)
			}

			return nil
		},
		WithPlanOptimizationLevel(OptimizationBasic),
		func(*planConfig) error {
			calls = append(calls, "last")

			return lastErr
		},
	})
	if !slices.Equal(calls, []string{"first", "middle", "last"}) || opts.level != compiler.Basic {
		t.Fatalf("option application: calls=%v level=%v", calls, opts.level)
	}

	if !errors.Is(err, firstErr) || !errors.Is(err, lastErr) {
		t.Fatalf("lost option failure: %v", err)
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok || !slices.Equal(joined.Unwrap(), []error{firstErr, lastErr}) {
		t.Fatalf("option failure order changed: %v", err)
	}
}

func TestPlanOptionsDefaultsAndValidation(t *testing.T) {
	for _, debug := range []bool{false, true} {
		initial := compiler.Basic
		if debug {
			initial = compiler.None
		}

		for _, setters := range [][]PlanOption{nil, {nil}} {
			opts, err := newPlanOptions(initial, debug, setters)
			if err != nil || opts.level != initial || opts.debug != debug {
				t.Fatalf("defaults changed: options=%+v err=%v", opts, err)
			}
		}

		for _, level := range []OptimizationLevel{OptimizationNone, OptimizationBasic, OptimizationFull, -1, 3, 99} {
			opts := planConfig{level: initial, debug: debug}
			err := WithPlanOptimizationLevel(level)(&opts)
			valid := level >= OptimizationNone && level <= OptimizationFull && (!debug || level == OptimizationNone)
			if valid {
				if err != nil || int(opts.level) != int(level) {
					t.Fatalf("debug=%v level=%v: options=%+v err=%v", debug, level, opts, err)
				}
			} else if err == nil || opts.level != initial {
				t.Fatalf("invalid level changed options: debug=%v level=%v options=%+v err=%v", debug, level, opts, err)
			}
		}
	}
}

func TestCompileDebugRejectsOptimizationBeforeHooks(t *testing.T) {
	var calls int
	engine := mustNewEngine(t, WithBeforeCompileHook(func(context.Context) error {
		calls++

		return nil
	}))
	t.Cleanup(func() { _ = engine.Close() })

	for _, level := range []OptimizationLevel{OptimizationBasic, OptimizationFull, -1, 3} {
		plan, err := engine.CompileDebug(t.Context(), NewAnonymousSource("RETURN 1"), WithPlanOptimizationLevel(level))
		if plan != nil || err == nil || calls != 0 {
			t.Fatalf("debug option validation: level=%v plan=%v err=%v hooks=%d", level, plan, err, calls)
		}
	}

	plan, err := engine.CompileDebug(t.Context(), NewAnonymousSource("RETURN 1"), WithPlanOptimizationLevel(OptimizationNone))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	if calls != 1 || plan.prog.Metadata.OptimizationLevel != int(OptimizationNone) {
		t.Fatalf("valid debug override: hooks=%d metadata=%+v", calls, plan.prog.Metadata)
	}
}
