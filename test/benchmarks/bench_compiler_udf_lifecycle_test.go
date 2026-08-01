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

const compilerUdfMemberStatementsQuery = `
FUNC read(value) (
  LET brand = value.product.brand
  VAR price = value["prices"]["current"]
  price = value.prices["sale"]
  value.metadata.lastSeen
  RETURN [brand, price]
)

RETURN read({
  product: { brand: "Ferret" },
  prices: { current: 1000, sale: 900 },
  metadata: { lastSeen: 1 }
})
`

const compilerUdfTransitiveCaptureQuery = `
LET base = 1

FUNC first(value) => second(value)
FUNC second(value) => third(value)
FUNC third(value) => base + value

RETURN first(1)
`

func BenchmarkCompilerCompileUdfLifecycle_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfLifecycleQuery, compiler.O0)
}

func BenchmarkCompilerCompileUdfLifecycle_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfLifecycleQuery, compiler.O1)
}

func BenchmarkCompilerCompileUdfMemberStatements_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfMemberStatementsQuery, compiler.O0)
}

func BenchmarkCompilerCompileUdfMemberStatements_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfMemberStatementsQuery, compiler.O1)
}

func BenchmarkCompilerCompileUdfTransitiveCapture_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfTransitiveCaptureQuery, compiler.O0)
}

func BenchmarkCompilerCompileUdfTransitiveCapture_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfTransitiveCaptureQuery, compiler.O1)
}
