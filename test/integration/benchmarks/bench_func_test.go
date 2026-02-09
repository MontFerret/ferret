package benchmarks_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

const (
	funcCallQuery = `
RETURN TEST(1,2,3,4,5,6)`
	func0CallQuery = `
RETURN TEST()`
	func1CallQuery = `
RETURN TEST(1)`
	func2CallQuery = `
RETURN TEST(1, 1)`
	func3CallQuery = `
RETURN TEST(1, 1, 1)`
	func4CallQuery = `
RETURN TEST(1, 1, 1, 1)`
)

func BenchmarkFunctionCall_O0(b *testing.B) {
	RunBenchmarkO0(b, funcCallQuery, vm.WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.True, nil
	}))
}

func BenchmarkFunctionCall_O1(b *testing.B) {
	RunBenchmarkO1(b, funcCallQuery, vm.WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.True, nil
	}))
}

func BenchmarkFunctionCall0_O0(b *testing.B) {
	RunBenchmarkO0(b, func0CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set0("TEST", func(ctx context.Context) (runtime.Value, error) {
			return runtime.String("test0"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall0_O1(b *testing.B) {
	RunBenchmarkO1(b, func0CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set0("TEST", func(ctx context.Context) (runtime.Value, error) {
			return runtime.String("test0"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall0Fallback_O0(b *testing.B) {
	RunBenchmarkO0(b, func0CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall0Fallback_O1(b *testing.B) {
	RunBenchmarkO1(b, func0CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall1_O0(b *testing.B) {
	RunBenchmarkO0(b, func1CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set1("TEST", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall1_O1(b *testing.B) {
	RunBenchmarkO1(b, func1CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set1("TEST", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall1Fallback_O0(b *testing.B) {
	RunBenchmarkO0(b, func1CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall1Fallback_O1(b *testing.B) {
	RunBenchmarkO1(b, func1CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall2_O0(b *testing.B) {
	RunBenchmarkO0(b, func2CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set2("TEST", func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall2_O1(b *testing.B) {
	RunBenchmarkO1(b, func2CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set2("TEST", func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall2Fallback_O0(b *testing.B) {
	RunBenchmarkO0(b, func2CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall2Fallback_O1(b *testing.B) {
	RunBenchmarkO1(b, func2CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall3_O0(b *testing.B) {
	RunBenchmarkO0(b, func3CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set3("TEST", func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall3_O1(b *testing.B) {
	RunBenchmarkO1(b, func3CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set3("TEST", func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall3Fallback_O0(b *testing.B) {
	RunBenchmarkO0(b, func3CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall3Fallback_O1(b *testing.B) {
	RunBenchmarkO1(b, func3CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall4_O0(b *testing.B) {
	RunBenchmarkO0(b, func4CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set4("TEST", func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall4_O1(b *testing.B) {
	RunBenchmarkO1(b, func4CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set4("TEST", func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall4Fallback_O0(b *testing.B) {
	RunBenchmarkO0(b, func4CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall4Fallback_O1(b *testing.B) {
	RunBenchmarkO1(b, func4CallQuery, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}
