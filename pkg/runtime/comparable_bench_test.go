package runtime_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var comparisonResult runtime.Ordering
var equalityResult runtime.Boolean
var hashResult uint64

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

func BenchmarkCompareValuesDuration(b *testing.B) {
	left := runtime.NewDuration(time.Second)
	right := runtime.NewDuration(2 * time.Second)
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

func BenchmarkEqualValuesDuration(b *testing.B) {
	left := runtime.NewDuration(time.Second)
	right := runtime.NewDuration(time.Second)
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

func BenchmarkCompareValuesMixedNumeric(b *testing.B) {
	ctx := context.Background()

	for _, test := range []struct {
		left  runtime.Value
		right runtime.Value
		name  string
	}{
		{name: "equal", left: runtime.NewInt(42), right: runtime.NewFloat(42)},
		{name: "fractional", left: runtime.NewInt(42), right: runtime.NewFloat(42.5)},
		{name: "large", left: runtime.NewInt64(1<<53 + 1), right: runtime.NewFloat(1 << 53)},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				result, err := runtime.CompareValues(ctx, test.left, test.right)
				if err != nil {
					b.Fatal(err)
				}
				comparisonResult = result
			}
		})
	}
}

func BenchmarkCompareValuesNaN(b *testing.B) {
	ctx := context.Background()
	nan := runtime.Value(runtime.NaN())

	b.ReportAllocs()

	for b.Loop() {
		result, err := runtime.CompareValues(ctx, nan, nan)
		if err != nil {
			b.Fatal(err)
		}
		comparisonResult = result
	}
}

func BenchmarkNumericHash(b *testing.B) {
	for _, test := range []struct {
		value runtime.Value
		name  string
	}{
		{name: "int", value: runtime.NewInt(42)},
		{name: "integral-float", value: runtime.NewFloat(42)},
		{name: "fractional-float", value: runtime.NewFloat(42.5)},
		{name: "nan", value: runtime.NewFloat(math.NaN())},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				hashResult = test.value.Hash()
			}
		})
	}
}
