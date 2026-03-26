package benchmarks_test

import "testing"

const (
	paramLoopShortQuery = `
FOR i IN 1..1000
  RETURN @test
`
)

func BenchmarkParamLoop_Short_O0(b *testing.B) {
	RunBenchmarkO0(b, paramLoopShortQuery, WithParam("test", "value"))
}

func BenchmarkParamLoop_Short_O1(b *testing.B) {
	RunBenchmarkO1(b, paramLoopShortQuery, WithParam("test", "value"))
}
