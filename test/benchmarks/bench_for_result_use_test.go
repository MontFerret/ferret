package benchmarks_test

import "testing"

const (
	discardedForResultQuery = `
FOR value IN 1..10000 {
  RETURN value * 2
}
`
	requiredForResultQuery = `
RETURN FOR value IN 1..10000 {
  RETURN value * 2
}
`
	discardedNestedForResultQuery = `
FOR outer IN 1..100 {
  RETURN FOR inner IN 1..100 {
    RETURN outer * inner
  }
}
`
	requiredNestedForResultQuery = `
RETURN FOR outer IN 1..100 {
  RETURN FOR inner IN 1..100 {
    RETURN outer * inner
  }
}
`
)

func BenchmarkForResultDiscarded_O0(b *testing.B) {
	RunBenchmarkO0(b, discardedForResultQuery)
}

func BenchmarkForResultDiscarded_O1(b *testing.B) {
	RunBenchmarkO1(b, discardedForResultQuery)
}

func BenchmarkForResultRequired_O0(b *testing.B) {
	RunBenchmarkO0(b, requiredForResultQuery)
}

func BenchmarkForResultRequired_O1(b *testing.B) {
	RunBenchmarkO1(b, requiredForResultQuery)
}

func BenchmarkForNestedResultDiscarded_O0(b *testing.B) {
	RunBenchmarkO0(b, discardedNestedForResultQuery)
}

func BenchmarkForNestedResultDiscarded_O1(b *testing.B) {
	RunBenchmarkO1(b, discardedNestedForResultQuery)
}

func BenchmarkForNestedResultRequired_O0(b *testing.B) {
	RunBenchmarkO0(b, requiredNestedForResultQuery)
}

func BenchmarkForNestedResultRequired_O1(b *testing.B) {
	RunBenchmarkO1(b, requiredNestedForResultQuery)
}
