package vm_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func hostInt(v runtime.Value) runtime.Int {
	i, ok := v.(runtime.Int)
	if !ok {
		return runtime.ZeroInt
	}

	return i
}

func hostSum(args ...runtime.Value) runtime.Int {
	sum := runtime.ZeroInt

	for _, arg := range args {
		sum += hostInt(arg)
	}

	return sum
}

func TestHostFunctionCall(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN TYPENAME(1)", "Int"),
		S("RETURN TYPENAME(1.1)", "Float"),
		S("WAIT(10) RETURN 1", 1),
		S("RETURN LENGTH([1,2,3])", 3),
		S("RETURN CONCAT('a', 'b', 'c')", "abc"),
		S("RETURN CONCAT(CONCAT('a', 'b'), 'c', CONCAT('d', 'e'))", "abcde", "Nested calls"),
		Array(`
		LET arr = []
		LET a = 1
		LET res = APPEND(arr, a)
		RETURN res
		`,
			[]any{1}, "Append to array"),
		S("LET duration = 10 WAIT(duration) RETURN 1", 1),
		Nil("RETURN (FALSE OR T::FAIL())?"),
		Nil("RETURN T::FAIL()?"),
		Array(`FOR i IN [1, 2, 3, 4]
				LET duration = 10
		
				WAIT(duration)
		
				RETURN i * 2`,
			[]any{2, 4, 6, 8}),

		S(`RETURN FIRST((FOR i IN 1..10 RETURN i * 2))`, 2),
		Array(`RETURN UNION((FOR i IN 0..5 RETURN i), (FOR i IN 6..10 RETURN i))`, []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
	})
}

func TestBuiltinFunctions(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN LENGTH([1,2,3])", 3),
		S("RETURN length([1,2,3])", 3),
		S("RETURN TYPENAME([1,2,3])", "Array"),
		S("RETURN TYPENAME({ a: 1, b: 2 })", "Object"),
		S("WAIT(10) RETURN 1", 1),
	})
}

func TestHostFunctionCallArities(t *testing.T) {
	builder := runtime.NewFunctionsBuilder()
	builder.A0().Add("FIX0", func(context.Context) (runtime.Value, error) {
		return runtime.NewString("fixed0"), nil
	})
	builder.A1().Add("FIX1", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
		return hostSum(arg), nil
	})
	builder.A2().Add("FIX2", func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
		return hostSum(arg1, arg2), nil
	})
	builder.A3().Add("FIX3", func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
		return hostSum(arg1, arg2, arg3), nil
	})
	builder.A4().Add("FIX4", func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
		return hostSum(arg1, arg2, arg3, arg4), nil
	})
	builder.Var().Add("VAR", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.NewInt(len(args)), nil
	})

	RunSpecs(t, []spec.Spec{
		S("RETURN FIX0()", "fixed0"),
		S("RETURN FIX1(1)", 1),
		S("RETURN FIX2(1, 2)", 3),
		S("RETURN FIX3(1, 2, 3)", 6),
		S("RETURN FIX4(1, 2, 3, 4)", 10),
		S("RETURN VAR()", 0),
		S("RETURN VAR(1)", 1),
		S("RETURN VAR(1, 2)", 2),
		S("RETURN VAR(1, 2, 3)", 3),
		S("RETURN VAR(1, 2, 3, 4)", 4),
		S("RETURN VAR(1, 2, 3, 4, 5)", 5),
		S("RETURN VAR(1, 2, 3, 4, 5, 6)", 6),
	}, vm.WithFunctionsBuilder(builder))
}

func TestHostFunctionProtectedCall(t *testing.T) {
	boom := errors.New("boom")

	RunSpecs(t, []spec.Spec{
		Nil("RETURN FAIL()?", "Protected host call should return none"),
		Nil("RETURN FAIL() ON ERROR RETURN NONE", "Explicit suppress should return none"),
		Nil("RETURN (FAIL() + 1) ON ERROR RETURN NONE", "Grouped explicit suppress should return none"),
		Error("RETURN FAIL()", "Non-protected host call should fail"),
		Error("RETURN FAIL() ON ERROR FAIL", "Explicit FAIL should preserve propagation"),
	}, vm.WithFunction("FAIL", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.None, boom
	}))
}

func TestHostFunctionRecoveryFallbackFailurePropagates(t *testing.T) {
	specs := []spec.Spec{
		spec.NewSpec("RETURN STEP() ON ERROR RETURN STEP()", "Fallback failure should escape instead of re-entering the same recovery tail").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{"boom-2"}},
		),
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		callCount := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			specs,
			vm.WithFunction("STEP", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				callCount++
				if callCount <= 2 {
					return runtime.None, fmt.Errorf("boom-%d", callCount)
				}

				return runtime.NewInt(99), nil
			}),
		)
	}
}

func TestHostFunctionLookupIsCaseSensitive(t *testing.T) {
	builder := runtime.NewFunctionsBuilder()
	builder.A0().Add("Foo", func(context.Context) (runtime.Value, error) {
		return runtime.NewString("upper"), nil
	})
	builder.A0().Add("foo", func(context.Context) (runtime.Value, error) {
		return runtime.NewString("lower"), nil
	})

	RunSpecs(t, []spec.Spec{
		Array("RETURN [Foo(), foo()]", []any{"upper", "lower"}),
		ErrorStr("RETURN FOO()", "unresolved function"),
	}, vm.WithFunctionsBuilder(builder))
}
