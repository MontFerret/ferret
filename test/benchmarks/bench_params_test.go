package benchmarks_test

import "testing"

const (
	paramLoopShortQuery = `
RETURN FOR i IN 1..1000
  RETURN @test
`

	paramLoopUDFQuery = `
FUNC read() => @test
RETURN FOR i IN 1..1000
  RETURN read()
`
)

func BenchmarkParamLoop_Short_None(b *testing.B) {
	RunBenchmarkNone(b, paramLoopShortQuery, WithParam("test", "value"))
}

func BenchmarkParamLoop_Short_Basic(b *testing.B) {
	RunBenchmarkBasic(b, paramLoopShortQuery, WithParam("test", "value"))
}

func BenchmarkParamLoop_Short_Full(b *testing.B) {
	RunBenchmarkFull(b, paramLoopShortQuery, WithParam("test", "value"))
}

func BenchmarkParamLoop_UDF_None(b *testing.B) {
	RunBenchmarkNone(b, paramLoopUDFQuery, WithParam("test", "value"))
}

func BenchmarkParamLoop_UDF_Basic(b *testing.B) {
	RunBenchmarkBasic(b, paramLoopUDFQuery, WithParam("test", "value"))
}

func BenchmarkParamLoop_UDF_Full(b *testing.B) {
	RunBenchmarkFull(b, paramLoopUDFQuery, WithParam("test", "value"))
}
