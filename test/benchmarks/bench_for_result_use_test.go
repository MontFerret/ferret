package benchmarks_test

import "testing"

const (
	discardedForResultQuery = `
FOR value IN 1..10000 {
  RETURN value * 2
}
`
	returnlessForResultQuery = `
FOR value IN 1..10000 {
  LET doubled = value * 2
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
	returnlessNestedForResultQuery = `
FOR outer IN 1..100 {
  FOR inner IN 1..100 {
    LET product = outer * inner
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

func BenchmarkForResultDiscarded_None(b *testing.B) {
	RunBenchmarkNone(b, discardedForResultQuery)
}

func BenchmarkForResultDiscarded_Basic(b *testing.B) {
	RunBenchmarkBasic(b, discardedForResultQuery)
}

func BenchmarkForResultDiscarded_Full(b *testing.B) {
	RunBenchmarkFull(b, discardedForResultQuery)
}

func BenchmarkForResultReturnless_None(b *testing.B) {
	RunBenchmarkNone(b, returnlessForResultQuery)
}

func BenchmarkForResultReturnless_Basic(b *testing.B) {
	RunBenchmarkBasic(b, returnlessForResultQuery)
}

func BenchmarkForResultReturnless_Full(b *testing.B) {
	RunBenchmarkFull(b, returnlessForResultQuery)
}

func BenchmarkForResultRequired_None(b *testing.B) {
	RunBenchmarkNone(b, requiredForResultQuery)
}

func BenchmarkForResultRequired_Basic(b *testing.B) {
	RunBenchmarkBasic(b, requiredForResultQuery)
}

func BenchmarkForResultRequired_Full(b *testing.B) {
	RunBenchmarkFull(b, requiredForResultQuery)
}

func BenchmarkForNestedResultDiscarded_None(b *testing.B) {
	RunBenchmarkNone(b, discardedNestedForResultQuery)
}

func BenchmarkForNestedResultDiscarded_Basic(b *testing.B) {
	RunBenchmarkBasic(b, discardedNestedForResultQuery)
}

func BenchmarkForNestedResultDiscarded_Full(b *testing.B) {
	RunBenchmarkFull(b, discardedNestedForResultQuery)
}

func BenchmarkForNestedResultReturnless_None(b *testing.B) {
	RunBenchmarkNone(b, returnlessNestedForResultQuery)
}

func BenchmarkForNestedResultReturnless_Basic(b *testing.B) {
	RunBenchmarkBasic(b, returnlessNestedForResultQuery)
}

func BenchmarkForNestedResultReturnless_Full(b *testing.B) {
	RunBenchmarkFull(b, returnlessNestedForResultQuery)
}

func BenchmarkForNestedResultRequired_None(b *testing.B) {
	RunBenchmarkNone(b, requiredNestedForResultQuery)
}

func BenchmarkForNestedResultRequired_Basic(b *testing.B) {
	RunBenchmarkBasic(b, requiredNestedForResultQuery)
}

func BenchmarkForNestedResultRequired_Full(b *testing.B) {
	RunBenchmarkFull(b, requiredNestedForResultQuery)
}
