package benchmarks_test

import "testing"

func BenchmarkLoop_Constants(b *testing.B) {
	RunBenchmarkNone(b, `
LET obj = { "a": 1 }
RETURN FOR i IN 1..100
  return obj.a
`)
}

func BenchmarkLoop_Constants_Full(b *testing.B) {
	RunBenchmarkFull(b, `
LET obj = { "a": 1 }
RETURN FOR i IN 1..100
  return obj.a
`)
}
