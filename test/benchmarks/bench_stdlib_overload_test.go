package benchmarks_test

import "testing"

const (
	trimDefaultQuery = `RETURN TRIM("  value  ")`
	trimCharsQuery   = `RETURN TRIM("xxvaluexx", "x")`
)

func BenchmarkStdlibTrimDefault_None(b *testing.B) {
	RunBenchmarkNone(b, trimDefaultQuery)
}

func BenchmarkStdlibTrimDefault_Full(b *testing.B) {
	RunBenchmarkFull(b, trimDefaultQuery)
}

func BenchmarkStdlibTrimChars_None(b *testing.B) {
	RunBenchmarkNone(b, trimCharsQuery)
}

func BenchmarkStdlibTrimChars_Full(b *testing.B) {
	RunBenchmarkFull(b, trimCharsQuery)
}
