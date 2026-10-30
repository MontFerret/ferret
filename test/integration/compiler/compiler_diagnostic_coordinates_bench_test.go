package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkDiagnosticCoordinates(b *testing.B) {
	const dialect = "RETURN QUERY `.x` IN @doc USING \"css\" /* ok */"
	const every = `RETURN WAITFOR EVENT "message" IN @obs TIMEOUT 5s EVERY 250ms`
	trailingComment := " /*" + strings.Repeat("a", 32<<10) + "*/"

	for _, prefix := range []struct{ name, text string }{{"ASCII", "// x\n"}, {"Unicode", "// 😀\n"}} {
		for _, test := range []struct{ name, query string }{
			{"While", "FOR 123 WHILE TRUE RETURN 1"},
			{"Query", dialect},
			{"Every", every},
			{"MissingIn", `RETURN WAITFOR EVENT "message" @source`},
			{"Multiple", "LET x =\nLET y =\nRETURN 1 +"},
			{"QueryTrailing32KiB", dialect + trailingComment},
			{"EveryTrailing32KiB", every + trailingComment},
		} {
			b.Run(prefix.name+"/"+test.name, func(b *testing.B) {
				c := mustNewCompiler(b)
				src := source.NewAnonymous(prefix.text + test.query)
				_, err := c.Compile(b.Context(), src)
				if test.name == "Multiple" {
					if set, ok := err.(*diagnostics.DiagnosticSet); !ok || set.Size() != 3 {
						b.Fatalf("expected three diagnostics, got %v", err)
					}
				} else if _, ok := err.(*diagnostics.Diagnostic); !ok {
					b.Fatalf("expected one diagnostic, got %v", err)
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, err := c.Compile(b.Context(), src); err == nil {
						b.Fatal("expected a diagnostic")
					}
				}
			})
		}
	}
}
