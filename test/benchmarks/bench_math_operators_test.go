package benchmarks_test

import "testing"

const mathOperatorsQuery = `
FOR i IN 1..1000
  LET sub = @left - @right
  LET mul = sub * @factor
  LET div = mul / @divisor
  RETURN div % @modulus
`

func BenchmarkMathOperators_O0(b *testing.B) {
	RunBenchmarkO0(
		b,
		mathOperatorsQuery,
		WithParam("left", 42),
		WithParam("right", 2),
		WithParam("factor", 3),
		WithParam("divisor", 4),
		WithParam("modulus", 7),
	)
}

func BenchmarkMathOperators_O1(b *testing.B) {
	RunBenchmarkO1(
		b,
		mathOperatorsQuery,
		WithParam("left", 42),
		WithParam("right", 2),
		WithParam("factor", 3),
		WithParam("divisor", 4),
		WithParam("modulus", 7),
	)
}
