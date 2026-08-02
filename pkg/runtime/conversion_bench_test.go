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

func BenchmarkToDateTimeRFC3339(b *testing.B) {
	ctx := context.Background()
	input := runtime.NewString("2026-08-02T12:00:00Z")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		conversionBenchmarkValue, conversionBenchmarkErr = runtime.ToDateTime(ctx, input)
	}
}

func BenchmarkToDateTimeEpochInt(b *testing.B) {
	ctx := context.Background()
	input := runtime.NewInt64(1_690_992_000)
	unit := runtime.NewString("s")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		conversionBenchmarkValue, conversionBenchmarkErr = runtime.ToDateTimeEpoch(ctx, input, unit)
	}
}

func BenchmarkToDateTimeEpochFloat(b *testing.B) {
	ctx := context.Background()
	input := runtime.NewFloat(1_690_992_000.5)
	unit := runtime.NewString("s")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		conversionBenchmarkValue, conversionBenchmarkErr = runtime.ToDateTimeEpoch(ctx, input, unit)
	}
}
