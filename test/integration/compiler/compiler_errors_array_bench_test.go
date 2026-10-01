package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkIncompleteArrayDiagnostic(b *testing.B) {
	for _, size := range []int{16, 4096} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			c := mustNewCompiler(b)
			src := source.NewAnonymous("RETURN [" + strings.Repeat("1, ", size-1) + "1")
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				if _, err := c.Compile(b.Context(), src); err == nil {
					b.Fatal("expected a syntax error")
				}
			}
		})
	}
}
