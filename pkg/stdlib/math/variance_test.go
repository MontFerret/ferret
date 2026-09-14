package math_test

import (
	"context"
	stdmath "math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestVarianceLargeOffset(t *testing.T) {
	for _, test := range []struct {
		call func(context.Context, runtime.Value) (runtime.Value, error)
		name string
		want runtime.Float
	}{
		{name: "population", call: math.PopulationVariance, want: 2},
		{name: "sample", call: math.SampleVariance, want: 2.5},
		{name: "stddev-population", call: math.StandardDeviationPopulation, want: runtime.Float(stdmath.Sqrt(2))},
		{name: "stddev-sample", call: math.StandardDeviationSample, want: runtime.Float(stdmath.Sqrt(2.5))},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := &aggregateList{values: []runtime.Value{
				runtime.Int(1_000_000_000_001), runtime.None,
				runtime.Float(1_000_000_000_002), runtime.String("2"),
				runtime.Int(1_000_000_000_003), runtime.True,
				runtime.Float(1_000_000_000_004), runtime.Int(1_000_000_000_005),
			}}
			got, err := test.call(t.Context(), source)
			if err != nil {
				t.Fatal(err)
			}

			assertMathResult(t, got, test.want)
			if source.calls != 1 {
				t.Fatalf("traversals = %d, want one", source.calls)
			}
		})
	}
}

func TestVarianceNonFiniteNumbers(t *testing.T) {
	for _, fn := range []func(context.Context, runtime.Value) (runtime.Value, error){
		math.PopulationVariance, math.SampleVariance,
		math.StandardDeviationPopulation, math.StandardDeviationSample,
	} {
		for _, value := range []runtime.Float{runtime.NaN(), runtime.Float(stdmath.Inf(1)), runtime.Float(stdmath.Inf(-1))} {
			for _, values := range [][]runtime.Value{{value}, {value, runtime.Int(1)}, {runtime.Int(1), value}} {
				got, err := fn(t.Context(), runtime.NewArrayWith(values...))
				if err != nil {
					t.Fatal(err)
				}

				assertMathResult(t, got, runtime.NaN())
			}
		}
	}
}
