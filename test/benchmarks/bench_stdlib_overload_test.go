package benchmarks_test

import "testing"

const (
	trimDefaultQuery = `RETURN TRIM("  value  ")`
	trimCharsQuery   = `RETURN TRIM("xxvaluexx", "x")`
)

func BenchmarkStdlibTrimDefault_O0(b *testing.B) {
	RunBenchmarkO0(b, trimDefaultQuery)
}

func BenchmarkStdlibTrimDefault_O1(b *testing.B) {
	RunBenchmarkO1(b, trimDefaultQuery)
}

func BenchmarkStdlibTrimChars_O0(b *testing.B) {
	RunBenchmarkO0(b, trimCharsQuery)
}

func BenchmarkStdlibTrimChars_O1(b *testing.B) {
	RunBenchmarkO1(b, trimCharsQuery)
}
