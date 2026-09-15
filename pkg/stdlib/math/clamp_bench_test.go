package math_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func BenchmarkClamp(b *testing.B) {
	for _, test := range []struct {
		min, max runtime.Value
		values   [3]runtime.Value
		name     string
	}{
		{name: "int", min: runtime.Int(2), max: runtime.Int(10), values: [3]runtime.Value{runtime.Int(-1), runtime.Int(5), runtime.Int(20)}},
		{name: "float", min: runtime.Float(0.5), max: runtime.Float(10.5), values: [3]runtime.Value{runtime.Float(-1.5), runtime.Float(5.5), runtime.Float(20.5)}},
		{name: "mixed", min: runtime.Float(0.5), max: runtime.Int(10), values: [3]runtime.Value{runtime.Int(-1), runtime.Float(5.5), runtime.Float(20.5)}},
	} {
		for index, name := range []string{"below", "inside", "above"} {
			b.Run(test.name+"/"+name, func(b *testing.B) {
				ctx := b.Context()
				value := test.values[index]
				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					result, err := math.Clamp(ctx, value, test.min, test.max)
					if err != nil {
						b.Fatal(err)
					}

					clampBenchmarkResult = result
				}
			})
		}
	}
}

var clampBenchmarkResult runtime.Value
