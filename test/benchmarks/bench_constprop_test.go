package benchmarks_test

import (
	"testing"
)

const constPropExpr = `
FOR i IN [1,2,3,4,5,6,7,8,9,10]
  LET v = (1 + 2) * (3 + 4) - (5 - 6)
  RETURN v
`

const constPropFloatDivisionByIntZeroExpr = `
FOR i IN 1..1000
  RETURN 1.0 / 0
`

func BenchmarkConstPropagation_O0(b *testing.B) {
	RunBenchmarkO0(b, constPropExpr)
}

func BenchmarkConstPropagation_O1(b *testing.B) {
	RunBenchmarkO1(b, constPropExpr)
}

func BenchmarkConstPropagationFloatDivisionByIntZero_O0(b *testing.B) {
	RunBenchmarkO0(b, constPropFloatDivisionByIntZeroExpr)
}

func BenchmarkConstPropagationFloatDivisionByIntZero_O1(b *testing.B) {
	RunBenchmarkO1(b, constPropFloatDivisionByIntZeroExpr)
}
