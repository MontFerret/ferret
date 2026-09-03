package benchmarks_test

import (
	"testing"
)

const sortQuery = `
LET strs = ["foo", "bar", "qaz", "abc"]

RETURN FOR s IN strs
	SORT s + "1"
	RETURN s
`

func BenchmarkForSort_None(b *testing.B) {
	RunBenchmarkNone(b, sortQuery)
}

func BenchmarkForSort_Basic(b *testing.B) {
	RunBenchmarkBasic(b, sortQuery)
}

func BenchmarkForSort_Full(b *testing.B) {
	RunBenchmarkFull(b, sortQuery)
}
