package benchmarks_test

import (
	"testing"
)

func BenchmarkArrayLiterals(b *testing.B) {
	RunBenchmark(b, `
RETURN ["foo", "bar", "qaz", "abc"]
`)
}

func BenchmarkObjectLiterals(b *testing.B) {
	RunBenchmark(b, `
RETURN { "foo": 1, "bar": 2, "qaz": 3, "abc": 4 }
`)
}
