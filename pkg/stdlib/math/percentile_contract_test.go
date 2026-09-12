package math_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestPercentileBoundaries(t *testing.T) {
	for _, test := range []struct {
		percent runtime.Value
		method  runtime.Value
		want    runtime.Value
		name    string
	}{
		{name: "rank-low", percent: runtime.Int(1), method: runtime.String("rank"), want: runtime.Int(1)},
		{name: "rank-high", percent: runtime.Int(100), method: runtime.String("rank"), want: runtime.Int(100)},
		{name: "rank-middle", percent: runtime.Int(50), method: runtime.String("rank"), want: runtime.Int(2)},
		{name: "rank-fractional-position", percent: runtime.Int(51), method: runtime.String("rank"), want: runtime.Float(3)},
		{name: "interpolation-low", percent: runtime.Int(1), method: runtime.String("interpolation"), want: runtime.Float(1.03)},
		{name: "interpolation-high", percent: runtime.Int(100), method: runtime.String("interpolation"), want: runtime.Int(100)},
		{name: "interpolation-middle", percent: runtime.Int(50), method: runtime.String("interpolation"), want: runtime.Float(2.5)},
		{name: "unknown-method", percent: runtime.Int(50), method: runtime.String("unknown"), want: runtime.Int(2)},
		{name: "method-case", percent: runtime.Int(50), method: runtime.String("INTERPOLATION"), want: runtime.Int(2)},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := runtime.NewArrayWith(runtime.Int(100), runtime.Float(3), runtime.Int(1), runtime.Int(2))
			before := source.String()
			got, err := math.Percentile(t.Context(), source, test.percent, test.method)
			if err != nil {
				t.Fatal(err)
			}

			assertMathResult(t, got, test.want)
			if source.String() != before {
				t.Fatalf("input changed: %s -> %s", before, source)
			}
		})
	}

	for _, method := range []runtime.String{"rank", "interpolation"} {
		got, err := math.Percentile(t.Context(), runtime.NewArrayWith(runtime.Int(9007199254740993)), runtime.Int(100), method)
		if err != nil {
			t.Fatal(err)
		}

		assertMathResult(t, got, runtime.Int(9007199254740993))
	}
}

func TestPercentileArgumentCompatibility(t *testing.T) {
	for _, percent := range []runtime.Value{runtime.Int(-1), runtime.Int(0), runtime.Int(101), runtime.Float(50), runtime.String("50"), runtime.None} {
		_, err := math.Percentile(t.Context(), runtime.NewArrayWith(runtime.Int(1)), percent)
		if err == nil {
			t.Fatalf("percentile %v (%T) accepted", percent, percent)
		}

		// The existing empty-list short-circuit deliberately precedes percentile validation.
		got, err := math.Percentile(t.Context(), runtime.NewArray(0), percent)
		if err != nil {
			t.Fatal(err)
		}

		assertMathResult(t, got, runtime.NaN())
	}

	for _, args := range [][]runtime.Value{
		nil,
		{runtime.NewArray(0)},
		{runtime.NewArray(0), runtime.Int(50), runtime.Int(1)},
		{runtime.NewArray(0), runtime.Int(50), runtime.String("rank"), runtime.Int(1)},
	} {
		if _, err := math.Percentile(t.Context(), args...); err == nil {
			t.Fatalf("invalid arguments accepted: %v", args)
		}
	}
}
