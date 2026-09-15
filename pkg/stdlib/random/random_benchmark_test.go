package random

import (
	"context"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkRandom(b *testing.B) {
	for name, call := range map[string]func(context.Context) (runtime.Value, error){
		"float": float0,
		"float_range": func(ctx context.Context) (runtime.Value, error) {
			return float2(ctx, runtime.Int(-10), runtime.Int(10))
		},
		"int": func(ctx context.Context) (runtime.Value, error) { return Int(ctx, runtime.Int(1), runtime.Int(6)) },
		"int_full": func(ctx context.Context) (runtime.Value, error) {
			return Int(ctx, runtime.Int(math.MinInt64), runtime.Int(math.MaxInt64))
		},
		"bool": Bool,
	} {
		b.Run(name, func(b *testing.B) {
			ctx := rnd.WithContext(b.Context(), rnd.NewSeed(42))
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				if _, err := call(ctx); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
