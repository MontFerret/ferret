package crypto_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/crypto"
)

func BenchmarkRandomToken(b *testing.B) {
	for _, size := range []runtime.Int{32, 65536} {
		b.Run(size.String(), func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				if _, err := crypto.RandomToken(context.Background(), size); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
