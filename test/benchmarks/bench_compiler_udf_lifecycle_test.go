package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

const compilerUdfLifecycleQuery = `
LET base = 1

FUNC outer(a) (
  VAR carried = base

  FUNC setCarried(v) (
    carried = v
    FUNC nested(c) => carried + a + c
    RETURN nested(1)
  )

  FUNC unusedInner() => carried

  RETURN setCarried(2)
)

FUNC unusedTop() => base

RETURN outer(3)
`

func BenchmarkCompilerCompileUdfLifecycle_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfLifecycleQuery, compiler.O0)
}

func BenchmarkCompilerCompileUdfLifecycle_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfLifecycleQuery, compiler.O1)
}
