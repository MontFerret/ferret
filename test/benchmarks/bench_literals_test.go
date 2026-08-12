package benchmarks_test

import (
	"testing"
)

const (
	arrayLiteralsQuery          = `RETURN ["foo", "bar", "qaz", "abc"]`
	arraySpreadQuery            = `LET left = ["foo", "bar"] LET right = ["qaz", "abc"] RETURN [...left, ...right]`
	objectLiteralsQuery         = `RETURN { "foo": 1, "bar": 2, "qaz": 3, "abc": 4 }`
	objectSpreadQuery           = `LET left = { "foo": 1, "bar": 2 } LET right = { "qaz": 3, "abc": 4 } RETURN { ...left, ...right }`
	objectComputedLiteralsQuery = `RETURN { ["foo"]: 1, ["bar"]: 2, [1]: 3 }`
)

func BenchmarkArrayLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, arrayLiteralsQuery)
}

func BenchmarkArrayLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, arrayLiteralsQuery)
}

func BenchmarkArraySpread_O0(b *testing.B) {
	RunBenchmarkO0(b, arraySpreadQuery)
}

func BenchmarkArraySpread_O1(b *testing.B) {
	RunBenchmarkO1(b, arraySpreadQuery)
}

func BenchmarkObjectLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, objectLiteralsQuery)
}

func BenchmarkObjectLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, objectLiteralsQuery)
}

func BenchmarkObjectSpread_O0(b *testing.B) {
	RunBenchmarkO0(b, objectSpreadQuery)
}

func BenchmarkObjectSpread_O1(b *testing.B) {
	RunBenchmarkO1(b, objectSpreadQuery)
}

func BenchmarkObjectComputedLiterals_O0(b *testing.B) {
	RunBenchmarkO0(b, objectComputedLiteralsQuery)
}

func BenchmarkObjectComputedLiterals_O1(b *testing.B) {
	RunBenchmarkO1(b, objectComputedLiteralsQuery)
}
