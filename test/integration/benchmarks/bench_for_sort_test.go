package benchmarks_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func BenchmarkForSort(b *testing.B) {
	RunBenchmark(b, `
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s
	RETURN s
`)
}
