package benchmarks_test

import (
	"testing"
)

const constPropExpr = `
RETURN FOR i IN [1,2,3,4,5,6,7,8,9,10]
  LET v = (1 + 2) * (3 + 4) - (5 - 6)
  RETURN v
`

func BenchmarkConstPropagation_None(b *testing.B) {
	RunBenchmarkNone(b, constPropExpr)
}

func BenchmarkConstPropagation_Basic(b *testing.B) {
	RunBenchmarkBasic(b, constPropExpr)
}

func BenchmarkConstPropagation_Full(b *testing.B) {
	RunBenchmarkFull(b, constPropExpr)
}
