package benchmarks_test

import (
	"testing"
)

func BenchmarkForSort(b *testing.B) {
	RunBenchmark(b, `
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s
	RETURN s
`)
}
