package arrays_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

func BenchmarkIntegerWidthArrayHints(b *testing.B) {
	input := runtime.NewArray(128)
	for index := range 128 {
		_ = input.Append(b.Context(), runtime.Int(index))
	}
	for _, name := range []string{"concat", "flatten"} {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var err error
				if name == "concat" {
					_, err = arrays.Concat(b.Context(), input, input)
				} else {
					_, err = arrays.Flatten(b.Context(), input)
				}
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
