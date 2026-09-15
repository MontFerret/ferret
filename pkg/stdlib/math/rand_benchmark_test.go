package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func BenchmarkLegacyRand(b *testing.B) {
	ctx := rnd.WithContext(context.Background(), rnd.NewSeed(0))
	b.ReportAllocs()

	for b.Loop() {
		if _, err := math.Rand(ctx); err != nil {
			b.Fatal(err)
		}
	}
}
