package benchmarks_test

import (
	"testing"
)

const (
	arrayLiteralsQuery          = `RETURN ["foo", "bar", "qaz", "abc"]`
	objectLiteralsQuery         = `RETURN { "foo": 1, "bar": 2, "qaz": 3, "abc": 4 }`
	objectComputedLiteralsQuery = `RETURN { ["foo"]: 1, ["bar"]: 2, [1]: 3 }`
)

func BenchmarkArrayLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, arrayLiteralsQuery)
}

func BenchmarkArrayLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, arrayLiteralsQuery)
}

func BenchmarkObjectLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, objectLiteralsQuery)
}

func BenchmarkObjectLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, objectLiteralsQuery)
}

func BenchmarkObjectComputedLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, objectComputedLiteralsQuery)
}

func BenchmarkObjectComputedLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, objectComputedLiteralsQuery)
}
