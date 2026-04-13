package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var benchmarkWarmupBindHostCallA1 = func(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return arg, nil
}

func BenchmarkWarmupBindHostCall_ArityMismatchSparseRegistry(b *testing.B) {
	descriptor := callDescriptor{
		DisplayName: "F",
		ArgCount:    2,
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		functions := mustBuildSparseBenchmarkFunctions(b)

		if _, err := warmupBindHostCall(descriptor, functions); !errors.Is(err, runtime.ErrInvalidArgumentNumber) {
			b.Fatalf("expected invalid argument number error, got %v", err)
		}
	}
}

func mustBuildSparseBenchmarkFunctions(b *testing.B) *runtime.Functions {
	b.Helper()

	builder := runtime.NewFunctionsBuilder()
	builder.A1().Add("F", benchmarkWarmupBindHostCallA1)

	functions, err := builder.Build()
	if err != nil {
		b.Fatalf("build functions: %v", err)
	}

	return functions
}
