package benchmarks_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
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

func withBuilder(add func(b *runtime.FunctionsBuilder)) vm.EnvironmentOption {
	builder := runtime.NewFunctionsBuilder()
	add(builder)
	return vm.WithFunctionsBuilder(builder)
}

func BenchmarkFunctionCall_None(b *testing.B) {
	RunBenchmarkNone(b, funcCallQuery, vm.WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.True, nil
	}))
}

func BenchmarkFunctionCall_Basic(b *testing.B) {
	RunBenchmarkBasic(b, funcCallQuery, vm.WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.True, nil
	}))
}

func BenchmarkFunctionCall_Full(b *testing.B) {
	RunBenchmarkFull(b, funcCallQuery, vm.WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.True, nil
	}))
}

func BenchmarkFunctionCall0_None(b *testing.B) {
	RunBenchmarkNone(b, func0CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A0().Add("TEST", func(ctx context.Context) (runtime.Value, error) {
			return runtime.String("test0"), nil
		})
	}))
}

func BenchmarkFunctionCall0_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func0CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A0().Add("TEST", func(ctx context.Context) (runtime.Value, error) {
			return runtime.String("test0"), nil
		})
	}))
}

func BenchmarkFunctionCall0_Full(b *testing.B) {
	RunBenchmarkFull(b, func0CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A0().Add("TEST", func(ctx context.Context) (runtime.Value, error) {
			return runtime.String("test0"), nil
		})
	}))
}

func BenchmarkFunctionCall0Fallback_None(b *testing.B) {
	RunBenchmarkNone(b, func0CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall0Fallback_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func0CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall0Fallback_Full(b *testing.B) {
	RunBenchmarkFull(b, func0CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall1_None(b *testing.B) {
	RunBenchmarkNone(b, func1CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A1().Add("TEST", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall1_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func1CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A1().Add("TEST", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall1_Full(b *testing.B) {
	RunBenchmarkFull(b, func1CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A1().Add("TEST", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall1Fallback_None(b *testing.B) {
	RunBenchmarkNone(b, func1CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall1Fallback_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func1CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall1Fallback_Full(b *testing.B) {
	RunBenchmarkFull(b, func1CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall2_None(b *testing.B) {
	RunBenchmarkNone(b, func2CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A2().Add("TEST", func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall2_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func2CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A2().Add("TEST", func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall2_Full(b *testing.B) {
	RunBenchmarkFull(b, func2CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A2().Add("TEST", func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall2Fallback_None(b *testing.B) {
	RunBenchmarkNone(b, func2CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall2Fallback_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func2CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall2Fallback_Full(b *testing.B) {
	RunBenchmarkFull(b, func2CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall3_None(b *testing.B) {
	RunBenchmarkNone(b, func3CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A3().Add("TEST", func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall3_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func3CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A3().Add("TEST", func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall3_Full(b *testing.B) {
	RunBenchmarkFull(b, func3CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A3().Add("TEST", func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall3Fallback_None(b *testing.B) {
	RunBenchmarkNone(b, func3CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall3Fallback_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func3CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall3Fallback_Full(b *testing.B) {
	RunBenchmarkFull(b, func3CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall4_None(b *testing.B) {
	RunBenchmarkNone(b, func4CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A4().Add("TEST", func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall4_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func4CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A4().Add("TEST", func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall4_Full(b *testing.B) {
	RunBenchmarkFull(b, func4CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.A4().Add("TEST", func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall4Fallback_None(b *testing.B) {
	RunBenchmarkNone(b, func4CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall4Fallback_Basic(b *testing.B) {
	RunBenchmarkBasic(b, func4CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}

func BenchmarkFunctionCall4Fallback_Full(b *testing.B) {
	RunBenchmarkFull(b, func4CallQuery, withBuilder(func(b *runtime.FunctionsBuilder) {
		b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		})
	}))
}
