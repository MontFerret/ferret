package benchmarks_test

import "testing"

func BenchmarkLoop_Constants(b *testing.B) {
	RunBenchmarkO0(b, `
LET obj = { "a": 1 }
FOR i IN 1..100
  return obj.a
`)
}

func BenchmarkLoop_Constants_O1(b *testing.B) {
	RunBenchmarkO1(b, `
LET obj = { "a": 1 }
FOR i IN 1..100
  return obj.a
`)
}
