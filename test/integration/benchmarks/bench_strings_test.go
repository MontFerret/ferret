package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
)

const addConstNumericQuery = `
FOR i IN 1..1000
  RETURN @base + 2
`

const addConstStringQuery = `
FOR i IN 1..1000
  RETURN @foo + " baz"
`

const templateLiteralSimpleQuery = "FOR i IN 1..1000 RETURN `Hello ${@name}!`"

const templateLiteralNumericQuery = "FOR i IN 1..1000 RETURN `sum=${@a + @b}`"

func BenchmarkAddConstNumeric_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstNumericQuery, vm.WithParam("base", 1))
}

func BenchmarkAddConstNumeric_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstNumericQuery, vm.WithParam("base", 1))
}

func BenchmarkAddConstString_O0(b *testing.B) {
	RunBenchmarkO0(b, addConstStringQuery, vm.WithParam("foo", "bar"))
}

func BenchmarkAddConstString_O1(b *testing.B) {
	RunBenchmarkO1(b, addConstStringQuery, vm.WithParam("foo", "bar"))
}

func BenchmarkTemplateLiteralSimple_O0(b *testing.B) {
	RunBenchmarkO0(b, templateLiteralSimpleQuery, vm.WithParam("name", "World"))
}

func BenchmarkTemplateLiteralSimple_O1(b *testing.B) {
	RunBenchmarkO1(b, templateLiteralSimpleQuery, vm.WithParam("name", "World"))
}

func BenchmarkTemplateLiteralNumeric_O0(b *testing.B) {
	RunBenchmarkO0(b, templateLiteralNumericQuery, vm.WithParam("a", 1), vm.WithParam("b", 2))
}

func BenchmarkTemplateLiteralNumeric_O1(b *testing.B) {
	RunBenchmarkO1(b, templateLiteralNumericQuery, vm.WithParam("a", 1), vm.WithParam("b", 2))
}
