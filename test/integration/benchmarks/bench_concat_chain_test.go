package benchmarks_test

import "testing"

const concatChainMixedQuery = `
FOR i IN 1..1000
  RETURN "a" + 1 + "b" + 2 + @name + "c" + 3 + @count + "d" + true + "e"
`

func BenchmarkConcatChainMixed_O0(b *testing.B) {
	RunBenchmarkO0(b, concatChainMixedQuery, WithParam("name", "X"), WithParam("count", 7))
}

func BenchmarkConcatChainMixed_O1(b *testing.B) {
	RunBenchmarkO1(b, concatChainMixedQuery, WithParam("name", "X"), WithParam("count", 7))
}
