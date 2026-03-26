package benchmarks_test

import (
	"testing"
)

const sortQuery = `
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s + "1"
	RETURN s
`

func BenchmarkForSort_O0(b *testing.B) {
	RunBenchmarkO0(b, sortQuery)
}

func BenchmarkForSort_O1(b *testing.B) {
	RunBenchmarkO1(b, sortQuery)
}
