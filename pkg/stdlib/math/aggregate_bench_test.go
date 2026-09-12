package math_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

var aggregateBenchmarkResult runtime.Value

func BenchmarkMathAggregates(b *testing.B) {
	functions := []struct {
		call func(context.Context, runtime.Value) (runtime.Value, error)
		name string
	}{
		{name: "sum", call: math.Sum},
		{name: "average", call: math.Average},
		{name: "min", call: math.Min},
		{name: "max", call: math.Max},
		{name: "variance", call: math.PopulationVariance},
		{name: "stddev", call: math.StandardDeviationSample},
		{name: "median", call: math.Median},
		{name: "percentile", call: func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return math.Percentile(ctx, value, runtime.Int(50), runtime.String("interpolation"))
		}},
	}

	for _, size := range []int{16, 1024} {
		values := make([]runtime.Value, size)
		for idx := range values {
			value := (idx * 37) % size
			if idx%2 == 0 {
				values[idx] = runtime.Int(value)
			} else {
				values[idx] = runtime.Float(value) + 0.5
			}
		}

		array := runtime.NewArrayWith(values...)
		for _, host := range []bool{false, true} {
			var source runtime.List = array
			if host {
				source = &struct{ runtime.List }{array}
			}

			for _, fn := range functions {
				b.Run(fmt.Sprintf("%s/%d/host=%t", fn.name, size, host), func(b *testing.B) {
					ctx := context.Background()
					b.ReportAllocs()
					b.ResetTimer()

					for b.Loop() {
						result, err := fn.call(ctx, source)
						if err != nil {
							b.Fatal(err)
						}

						aggregateBenchmarkResult = result
					}
				})
			}
		}
	}
}
