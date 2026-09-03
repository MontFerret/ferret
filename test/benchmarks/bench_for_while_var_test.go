package benchmarks_test

import "testing"

const forWhileVarQuery = `
VAR i = 0
RETURN FOR WHILE i < 100
  LET current = i
  i = i + 1
  RETURN current
`

func BenchmarkForWhileVar_None(b *testing.B) {
	RunBenchmarkNone(b, forWhileVarQuery)
}

func BenchmarkForWhileVar_Basic(b *testing.B) {
	RunBenchmarkBasic(b, forWhileVarQuery)
}

func BenchmarkForWhileVar_Full(b *testing.B) {
	RunBenchmarkFull(b, forWhileVarQuery)
}
