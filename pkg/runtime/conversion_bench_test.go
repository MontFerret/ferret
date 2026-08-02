package runtime_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	conversionBenchmarkValue runtime.Value
	conversionBenchmarkErr   error
)

func BenchmarkToIntString(b *testing.B) {
	ctx := context.Background()
	input := runtime.NewString("12345")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		conversionBenchmarkValue, conversionBenchmarkErr = runtime.ToInt(ctx, input)
	}
}

func BenchmarkToIntInvalidString(b *testing.B) {
	ctx := context.Background()
	input := runtime.NewString("invalid")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		conversionBenchmarkValue, conversionBenchmarkErr = runtime.ToInt(ctx, input)
	}
}

func BenchmarkToFloatString(b *testing.B) {
	ctx := context.Background()
	input := runtime.NewString("12345.5")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		conversionBenchmarkValue, conversionBenchmarkErr = runtime.ToFloat(ctx, input)
	}
}

func BenchmarkToFloatInvalidString(b *testing.B) {
	ctx := context.Background()
	input := runtime.NewString("invalid")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		conversionBenchmarkValue, conversionBenchmarkErr = runtime.ToFloat(ctx, input)
	}
}
