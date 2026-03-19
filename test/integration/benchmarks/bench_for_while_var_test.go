package benchmarks_test

import "testing"

const forWhileVarQuery = `
VAR i = 0
FOR WHILE i < 100
  LET current = i
  i = i + 1
  RETURN current
`

func BenchmarkForWhileVar_O0(b *testing.B) {
	RunBenchmarkO0(b, forWhileVarQuery)
}

func BenchmarkForWhileVar_O1(b *testing.B) {
	RunBenchmarkO1(b, forWhileVarQuery)
}
