package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkIntegerWidthVMArrayHints(b *testing.B) {
	input := runtime.NewArray(128)
	for index := range 128 {
		_ = input.Append(b.Context(), runtime.Int(index))
	}

	for _, name := range []string{"flatten", "distinct"} {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var err error
				if name == "flatten" {
					_, err = arrayFlatten(b.Context(), input, 1)
				} else {
					_, err = arrayDistinct(b.Context(), input)
				}

				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
