package runtime_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkArraySlice(b *testing.B) {
	for _, size := range []int{0, 16, 1024} {
		b.Run(runtime.Int(size).String(), func(b *testing.B) {
			input := runtime.NewSizedArray(size)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				if _, err := input.Slice(b.Context(), 0, runtime.Int(size)); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
