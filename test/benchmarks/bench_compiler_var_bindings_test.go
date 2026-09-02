package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

const compilerVarBindingsQuery = `
VAR total = 1
total += 2

VAR carried = 0

FUNC bump(v) {
  carried = v
  RETURN carried
}

LET ignored = (
  FOR item IN [1, 2, 3]
    carried = item
    LET _ = bump(item + total)
    RETURN item
)

RETURN carried
`

func BenchmarkCompilerCompileVarBindings_None(b *testing.B) {
	benchmarkCompileQuery(b, compilerVarBindingsQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileVarBindings_Basic(b *testing.B) {
	benchmarkCompileQuery(b, compilerVarBindingsQuery, compiler.OptimizationBasic)
}

func BenchmarkCompilerCompileVarBindings_Full(b *testing.B) {
	benchmarkCompileQuery(b, compilerVarBindingsQuery, compiler.OptimizationFull)
}
