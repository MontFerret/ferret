package benchmarks_test

import (
	"testing"
)

const (
	addConstNumericQuery = `
LET base = 1
FOR i IN 1..1000
  RETURN base + 2
`

	addConstNumericWithParamQuery = `
FOR i IN 1..1000
  RETURN @base + 2
`

	addConstStringQuery = `
LET foo = "bar"
FOR i IN 1..1000
  RETURN foo + " baz"`

	addConstStringWithParamQuery = `
FOR i IN 1..1000
  RETURN @foo + " baz"
`

	templateLiteralSimpleQuery = "FOR i IN 1..1000 RETURN `Hello ${@name}!`"

	templateLiteralNumericQuery = "FOR i IN 1..1000 RETURN `sum=${@a + @b}`"
)

func BenchmarkAddNumeric_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstNumericQuery)
}

func BenchmarkAddNumeric_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstNumericQuery)
}

func BenchmarkAddConstNumericWithParam_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstNumericWithParamQuery, WithParam("base", 1))
}

func BenchmarkAddConstNumericWithParam_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstNumericWithParamQuery, WithParam("base", 1))
}

func BenchmarkAddConstString_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstStringQuery)
}

func BenchmarkAddConstString_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstStringQuery)
}

func BenchmarkAddConstStringWithParam_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstStringWithParamQuery, WithParam("foo", "bar"))
}

func BenchmarkAddConstStringWithParam_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstStringWithParamQuery, WithParam("foo", "bar"))
}

func BenchmarkTemplateLiteralSimple_O0(b *testing.B) {
	RunBenchmarkO0(b, templateLiteralSimpleQuery, WithParam("name", "World"))
}

func BenchmarkTemplateLiteralSimple_O1(b *testing.B) {
	RunBenchmarkO1(b, templateLiteralSimpleQuery, WithParam("name", "World"))
}

func BenchmarkTemplateLiteralNumeric_O0(b *testing.B) {
	RunBenchmarkO0(b, templateLiteralNumericQuery, WithParam("a", 1), WithParam("b", 2))
}

func BenchmarkTemplateLiteralNumeric_O1(b *testing.B) {
	RunBenchmarkO1(b, templateLiteralNumericQuery, WithParam("a", 1), WithParam("b", 2))
}
