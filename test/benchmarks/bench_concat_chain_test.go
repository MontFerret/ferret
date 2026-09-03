package benchmarks_test

import "testing"

const concatChainMixedQuery = `
RETURN FOR i IN 1..1000
  RETURN "a" + 1 + "b" + 2 + @name + "c" + 3 + @count + "d" + true + "e"
`

func BenchmarkConcatChainMixed_None(b *testing.B) {
	RunBenchmarkNone(b, concatChainMixedQuery, WithParam("name", "X"), WithParam("count", 7))
}

func BenchmarkConcatChainMixed_Basic(b *testing.B) {
	RunBenchmarkBasic(b, concatChainMixedQuery, WithParam("name", "X"), WithParam("count", 7))
}

func BenchmarkConcatChainMixed_Full(b *testing.B) {
	RunBenchmarkFull(b, concatChainMixedQuery, WithParam("name", "X"), WithParam("count", 7))
}
