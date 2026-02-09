package benchmarks_test

import (
	"testing"
)

const constPropExpr = `
FOR i IN [1,2,3,4,5,6,7,8,9,10]
  LET v = (1 + 2) * (3 + 4) - (5 - 6)
  RETURN v
`

func BenchmarkConstPropagation_O0(b *testing.B) {
	RunBenchmarkO0(b, constPropExpr)
}

func BenchmarkConstPropagation_O1(b *testing.B) {
	RunBenchmarkO1(b, constPropExpr)
}
