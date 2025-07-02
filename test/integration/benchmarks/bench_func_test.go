package benchmarks_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func BenchmarkFunctionCall(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1,2,3,4,5,6)
`, vm.WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.True, nil
	}))
}

func BenchmarkFunctionCall0(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST()
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set0("TEST", func(ctx context.Context) (runtime.Value, error) {
			return runtime.String("test0"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall0Fallback(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST()
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall1(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1)
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set1("TEST", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall1Fallback(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1)
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall2(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1, 1)
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set2("TEST", func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall2Fallback(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1, 1)
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall3(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1, 1, 1)
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set3("TEST", func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall3Fallback(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1, 1, 1)
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall4(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1, 1, 1, 1)
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set4("TEST", func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}

func BenchmarkFunctionCall4Fallback(b *testing.B) {
	RunBenchmark(b, `
	RETURN TEST(1, 1, 1, 1)
`, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Build()))
}
