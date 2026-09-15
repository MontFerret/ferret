package random

import (
	"context"
	"fmt"
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

func BenchmarkRandomCollections(b *testing.B) {
	for name, call := range collectionOperations {
		for _, size := range []int{10, 1000} {
			b.Run(fmt.Sprintf("%s/Array/%d", name, size), func(b *testing.B) {
				values := make([]runtime.Value, size)
				for i := range values {
					values[i] = runtime.Int(i)
				}

				benchmarkRandomCollection(b, call, runtime.NewArrayOf(values))
			})
		}
	}

	b.Run("choice/traversal-only/1000", func(b *testing.B) {
		values := make([]runtime.Value, 1000)
		for i := range values {
			values[i] = runtime.Int(i)
		}

		benchmarkRandomCollection(b, Choice, &traversalValues{values: values})
	})
}

func benchmarkRandomCollection(b *testing.B, call func(context.Context, runtime.Value) (runtime.Value, error), values runtime.List) {
	b.Helper()

	ctx := rnd.WithContext(b.Context(), rnd.NewSeed(42))
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := call(ctx, values); err != nil {
			b.Fatal(err)
		}
	}
}
