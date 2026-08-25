package ferret

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func TestCompileAppliesUniversalOptimizationLevelsPerCompilation(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t, WithCompilerOptions(compiler.WithOptimizationLevel(compiler.O1)))
	t.Cleanup(func() { _ = engine.Close() })

	levels := []OptimizationLevel{
		OptimizationNone,
		OptimizationBasic,
		OptimizationFull,
		OptimizationAggressive,
	}

	for _, level := range levels {
		level := level
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			t.Parallel()

			compiled, err := engine.Compile(
				context.Background(),
				NewAnonymousSource("LET value = 1 + 2 RETURN value"),
				WithOptimizationLevel(level),
			)
			if err != nil {
				t.Fatalf("compile O%d: %v", level, err)
			}
			defer compiled.Close()

			plan, ok := compiled.(*Plan)
			if !ok {
				t.Fatalf("unexpected plan type %T", compiled)
			}

			if got, want := plan.prog.Metadata.OptimizationLevel, int(level); got != want {
				t.Fatalf("optimization level = %d, want %d", got, want)
			}
		})
	}

	compiled, err := engine.Compile(context.Background(), NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatalf("compile with engine default: %v", err)
	}
	defer compiled.Close()

	plan, ok := compiled.(*Plan)
	if !ok {
		t.Fatalf("unexpected plan type %T", compiled)
	}

	if got := plan.prog.Metadata.OptimizationLevel; got != int(compiler.O1) {
		t.Fatalf("engine compiler mutated to O%d, want O%d", got, compiler.O1)
	}
}

func TestCompileDebugUsesEffectiveO0WithUniversalOverride(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t)
	t.Cleanup(func() { _ = engine.Close() })

	compiled, err := engine.CompileDebug(
		context.Background(),
		NewAnonymousSource("LET value = 1 RETURN value"),
		WithOptimizationLevel(OptimizationAggressive),
	)
	if err != nil {
		t.Fatalf("compile debug: %v", err)
	}
	defer compiled.Close()

	plan, ok := compiled.(*Plan)
	if !ok {
		t.Fatalf("unexpected plan type %T", compiled)
	}

	if got := plan.prog.Metadata.OptimizationLevel; got != int(compiler.O0) {
		t.Fatalf("debug optimization level = O%d, want O%d", got, compiler.O0)
	}

	if len(plan.prog.Metadata.DebugPoints) == 0 {
		t.Fatal("debug compilation did not emit debug points")
	}
}

func TestCompileAppliesPlanOptionsInOrderAndJoinsFailures(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t)
	t.Cleanup(func() { _ = engine.Close() })

	firstErr := errors.New("first plan option failed")
	secondErr := errors.New("second plan option failed")
	var calls []string

	plan, err := engine.Compile(
		context.Background(),
		NewAnonymousSource("RETURN 1"),
		func(PlanOptions) error {
			calls = append(calls, "first")

			return firstErr
		},
		nil,
		func(PlanOptions) error {
			calls = append(calls, "middle")

			return nil
		},
		func(PlanOptions) error {
			calls = append(calls, "second")

			return secondErr
		},
	)
	if plan != nil {
		_ = plan.Close()

		t.Fatalf("plan = %T, want nil", plan)
	}

	if !slices.Equal(calls, []string{"first", "middle", "second"}) {
		t.Fatalf("plan option calls = %v", calls)
	}

	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("compile error = %v, want both option errors", err)
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("compile error = %T, want joined error", err)
	}

	if failures := joined.Unwrap(); !slices.Equal(failures, []error{firstErr, secondErr}) {
		t.Fatalf("plan option failures = %v, want ordered sentinels", failures)
	}
}
