package sdk_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

func BenchmarkIntegerWidthToSlice(b *testing.B) {
	input := runtime.NewArray(128)
	for index := range 128 {
		_ = input.Append(b.Context(), runtime.Int(index))
	}
	mapper := func(_ context.Context, value, _ runtime.Value) (int64, error) { return int64(value.(runtime.Int)), nil }
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := sdk.ToSlice(b.Context(), input, mapper); err != nil {
			b.Fatal(err)
		}
	}
}
