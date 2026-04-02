package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

const compilerVarBindingsQuery = `
VAR total = 1
total += 2

VAR carried = 0

FUNC bump(v) (
  carried = v
  RETURN carried
)

LET ignored = (
  FOR item IN [1, 2, 3]
    carried = item
    LET _ = bump(item + total)
    RETURN item
)

RETURN carried
`

func BenchmarkCompilerCompileVarBindings_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerVarBindingsQuery, compiler.O0)
}

func BenchmarkCompilerCompileVarBindings_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerVarBindingsQuery, compiler.O1)
}
