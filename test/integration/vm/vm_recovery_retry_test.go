package vm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	. "github.com/MontFerret/ferret/v2/test/spec/mock"
)

func TestHostFunctionRetryRecoverySuccess(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		callCount := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				S("RETURN STEP() ON ERROR RETRY 2 DELAY 0 BACKOFF EXPONENTIAL", 99, "Retry should return the first successful attempt"),
			},
			vm.WithFunction("STEP", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				callCount++
				if callCount <= 2 {
					return runtime.None, fmt.Errorf("boom-%d", callCount)
				}

				return runtime.NewInt(99), nil
			}),
		)

		if got, want := callCount, 3; got != want {
			t.Fatalf("unexpected retry attempt count for O%d: got %d, want %d", level, got, want)
		}
	}
}

func TestHostFunctionRetryExhaustionPropagatesFinalError(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		callCount := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				spec.NewSpec("RETURN STEP() ON ERROR RETRY 1", "Retry without fallback should propagate the final unhandled failure").Expect().ExecError(
					ShouldBeRuntimeError,
					&ExpectedRuntimeError{Contains: []string{"boom-2"}},
				),
			},
			vm.WithFunction("STEP", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				callCount++
				return runtime.None, fmt.Errorf("boom-%d", callCount)
			}),
		)

		if got, want := callCount, 2; got != want {
			t.Fatalf("unexpected propagated retry attempt count for O%d: got %d, want %d", level, got, want)
		}
	}
}

func TestHostFunctionRetryFallbackFailurePropagates(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		callCount := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				spec.NewSpec("RETURN STEP() ON ERROR RETRY 1 OR RETURN STEP()", "Retry fallback failure should escape instead of re-entering retry").Expect().ExecError(
					ShouldBeRuntimeError,
					&ExpectedRuntimeError{Contains: []string{"boom-3"}},
				),
			},
			vm.WithFunction("STEP", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				callCount++
				return runtime.None, fmt.Errorf("boom-%d", callCount)
			}),
		)

		if got, want := callCount, 3; got != want {
			t.Fatalf("unexpected retry fallback attempt count for O%d: got %d, want %d", level, got, want)
		}
	}
}

func TestGroupedForRetryRecovery(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		callCount := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				Array(`LET xs = (FOR i IN [1, 2] LET y = STEP() RETURN y + i) ON ERROR RETRY 1 OR RETURN []
RETURN xs`, []any{11, 12}, "Grouped FOR retry should restart after cleanup without leaking partial results"),
			},
			vm.WithFunction("STEP", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				callCount++
				if callCount == 1 {
					return runtime.None, fmt.Errorf("boom-%d", callCount)
				}

				return runtime.NewInt(10), nil
			}),
		)

		if got, want := callCount, 3; got != want {
			t.Fatalf("unexpected grouped FOR retry attempt count for O%d: got %d, want %d", level, got, want)
		}
	}
}

func TestWaitForPredicateRetryRecovery(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		callCount := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				S(`LET token = WAITFOR VALUE STEP() TIMEOUT 20ms EVERY 0 ON TIMEOUT RETURN "timeout" ON ERROR RETRY 2 DELAY 0 OR RETURN "error"
RETURN token`, "ok", "WAITFOR predicate should retry runtime failures and keep timeout handling separate"),
			},
			vm.WithFunction("STEP", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				callCount++
				if callCount <= 2 {
					return runtime.None, fmt.Errorf("boom-%d", callCount)
				}

				return runtime.NewString("ok"), nil
			}),
		)

		if got, want := callCount, 3; got != want {
			t.Fatalf("unexpected WAITFOR predicate retry attempt count for O%d: got %d, want %d", level, got, want)
		}
	}
}

func TestWaitForEventRetryRecovery(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		successCalls := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				NotNil(`RETURN WAITFOR EVENT "test" IN SOURCE() TIMEOUT 20ms ON TIMEOUT RETURN "timeout" ON ERROR RETRY 2 DELAY 0 OR RETURN "error"`, "WAITFOR EVENT should retry source evaluation failures"),
			},
			vm.WithFunction("SOURCE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				successCalls++
				if successCalls <= 2 {
					return runtime.None, fmt.Errorf("boom-%d", successCalls)
				}

				return NewObservable([]runtime.Value{NewTestEventType("match")}), nil
			}),
		)

		if got, want := successCalls, 3; got != want {
			t.Fatalf("unexpected WAITFOR EVENT retry attempt count for O%d: got %d, want %d", level, got, want)
		}

		timeoutCalls := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				S(`RETURN WAITFOR EVENT "test" IN SOURCE() TIMEOUT 1ms ON TIMEOUT RETURN "timeout" ON ERROR RETRY 2 DELAY 0 OR RETURN "error"`, "timeout", "WAITFOR EVENT timeout should not be retried by ON ERROR RETRY"),
			},
			vm.WithFunction("SOURCE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				timeoutCalls++
				return NewBlockingObservable(), nil
			}),
		)

		if got, want := timeoutCalls, 1; got != want {
			t.Fatalf("unexpected WAITFOR EVENT timeout source call count for O%d: got %d, want %d", level, got, want)
		}
	}
}
