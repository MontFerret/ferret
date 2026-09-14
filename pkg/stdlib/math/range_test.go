package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestRangeCompatibility(t *testing.T) {
	for _, args := range [][]runtime.Value{
		{runtime.Int(1), runtime.Int(4)},
		{runtime.Float(-0.75), runtime.Float(1.1), runtime.Float(0.5)},
		{runtime.Int(4), runtime.Int(1), runtime.Int(-1)},
		{runtime.Int(1), runtime.Int(4), runtime.Int(0)},
		{runtime.Int(1)},
	} {
		got, err := math.Range(t.Context(), args...)
		want, wantErr := arrays.Range(t.Context(), args...)
		if (err == nil) != (wantErr == nil) || got.String() != want.String() {
			t.Fatalf("range compatibility: %v, %v; want %v, %v", got, err, want, wantErr)
		}
	}
}

var rangeBenchmarkResult runtime.Value

func BenchmarkRange(b *testing.B) {
	ctx := context.Background()
	tests := []struct {
		name string
		args []runtime.Value
	}{
		{name: "default_step", args: []runtime.Value{runtime.NewInt(1), runtime.NewInt(100)}},
		{name: "custom_step", args: []runtime.Value{runtime.NewInt(10), runtime.NewInt(100), runtime.NewInt(2)}},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()

			for range b.N {
				result, err := math.Range(ctx, test.args...)
				if err != nil {
					b.Fatal(err)
				}

				rangeBenchmarkResult = result
			}
		})
	}
}
