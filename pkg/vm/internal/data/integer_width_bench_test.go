package data

import (
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkIntegerWidthObjectComparison(b *testing.B) {
	values := make(map[string]runtime.Value, 128)
	for index := range 128 {
		values[strconv.Itoa(index)] = runtime.Int(index)
	}

	left := runtime.NewObjectWith(values)
	right := runtime.NewObjectWith(values)
	b.ReportAllocs()
	for b.Loop() {
		if equal, err := equalObjectLike(b.Context(), left, right); err != nil || !equal {
			b.Fatalf("comparison = %v, %v", equal, err)
		}
	}
}
