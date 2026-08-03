package runtime_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var comparisonResult runtime.Ordering
var equalityResult runtime.Boolean

func BenchmarkCompareValuesNative(b *testing.B) {
	left := runtime.NewInt(41)
	right := runtime.NewInt(42)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := runtime.CompareValues(ctx, left, right)
		if err != nil {
			b.Fatal(err)
		}
		comparisonResult = result
	}
}

func BenchmarkCompareValuesNested(b *testing.B) {
	left := runtime.NewArrayWith(
		runtime.NewInt(1),
		runtime.NewObjectWith(map[string]runtime.Value{
			"items": runtime.NewArrayWith(
				runtime.NewString("alpha"),
				runtime.NewFloat(2.5),
			),
		}),
	)
	right := runtime.NewArrayWith(
		runtime.NewInt(1),
		runtime.NewObjectWith(map[string]runtime.Value{
			"items": runtime.NewArrayWith(
				runtime.NewString("alpha"),
				runtime.NewFloat(3.5),
			),
		}),
	)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := runtime.CompareValues(ctx, left, right)
		if err != nil {
			b.Fatal(err)
		}
		comparisonResult = result
	}
}

func BenchmarkCompareValuesHost(b *testing.B) {
	left := comparisonBenchmarkValue{value: 41}
	right := comparisonBenchmarkValue{value: 42}
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := runtime.CompareValues(ctx, left, right)
		if err != nil {
			b.Fatal(err)
		}
		comparisonResult = result
	}
}

func BenchmarkEqualValuesHost(b *testing.B) {
	left := comparisonBenchmarkValue{value: 42}
	right := comparisonBenchmarkValue{value: 42}
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := runtime.EqualValues(ctx, left, right)
		if err != nil {
			b.Fatal(err)
		}
		equalityResult = result
	}
}
