package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

const compilerUdfLifecycleQuery = `
LET base = 1

FUNC outer(a) {
  VAR carried = base

  FUNC setCarried(v) {
    carried = v
    FUNC nested(c) => carried + a + c
    RETURN nested(1)
  }

  FUNC unusedInner() => carried

  RETURN setCarried(2)
}

FUNC unusedTop() => base

RETURN outer(3)
`

const compilerUdfMemberStatementsQuery = `
FUNC read(value) {
  LET brand = value.product.brand
  VAR price = value["prices"]["current"]
  price = value.prices["sale"]
  value.metadata.lastSeen
  RETURN [brand, price]
}

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
	benchmarkCompileQuery(b, compilerUdfLifecycleQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileUdfLifecycle_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfLifecycleQuery, compiler.OptimizationFull)
}

func BenchmarkCompilerCompileUdfMemberStatements_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfMemberStatementsQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileUdfMemberStatements_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfMemberStatementsQuery, compiler.OptimizationFull)
}

func BenchmarkCompilerCompileUdfTransitiveCapture_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfTransitiveCaptureQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileUdfTransitiveCapture_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfTransitiveCaptureQuery, compiler.OptimizationFull)
}

func BenchmarkCompilerCompileUdfTransitiveCapture_Basic(b *testing.B) {
	benchmarkCompileQuery(b, compilerUdfTransitiveCaptureQuery, compiler.OptimizationBasic)
}
