package objects

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkIntegerWidthListCopy(b *testing.B) {
	input := runtime.NewArray(128)
	for index := range 128 {
		_ = input.Append(b.Context(), runtime.Int(index))
	}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := copyListValues(b.Context(), input); err != nil {
			b.Fatal(err)
		}
	}
}
