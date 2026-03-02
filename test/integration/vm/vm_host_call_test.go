package vm_test

import (
	"context"
	"errors"
	"testing"

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

	RunUseCases(t, []UseCase{
		Case("RETURN FIX0()", "fixed0"),
		Case("RETURN FIX1(1)", 1),
		Case("RETURN FIX2(1, 2)", 3),
		Case("RETURN FIX3(1, 2, 3)", 6),
		Case("RETURN FIX4(1, 2, 3, 4)", 10),
		Case("RETURN VAR()", 0),
		Case("RETURN VAR(1)", 1),
		Case("RETURN VAR(1, 2)", 2),
		Case("RETURN VAR(1, 2, 3)", 3),
		Case("RETURN VAR(1, 2, 3, 4)", 4),
		Case("RETURN VAR(1, 2, 3, 4, 5)", 5),
		Case("RETURN VAR(1, 2, 3, 4, 5, 6)", 6),
	}, vm.WithFunctionsBuilder(builder))
}

func TestHostFunctionProtectedCall(t *testing.T) {
	boom := errors.New("boom")

	RunUseCases(t, []UseCase{
		CaseNil("RETURN FAIL()?", "Protected host call should return none"),
		CaseRuntimeError("RETURN FAIL()", "Non-protected host call should fail"),
	}, vm.WithFunction("FAIL", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.None, boom
	}))
}
