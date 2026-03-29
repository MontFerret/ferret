package benchmarks_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkCompilerCompileConcatChain_O1(b *testing.B) {
	benchmarkCompileQuery(b, buildConcatCompileQuery(false), compiler.O1)
}

func BenchmarkCompilerCompileStringAppend_O1(b *testing.B) {
	benchmarkCompileQuery(b, buildConcatCompileQuery(true), compiler.O1)
}

func benchmarkCompileQuery(b *testing.B, query string, level compiler.OptimizationLevel) {
	b.Helper()

	compilerInstance := compiler.New(compiler.WithOptimizationLevel(level))
	src := source.New("concat_benchmark", query)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := compilerInstance.Compile(src); err != nil {
			b.Fatalf("compile failed: %v", err)
		}
	}
}

func buildConcatCompileQuery(appendStyle bool) string {
	var b strings.Builder

	if appendStyle {
		b.WriteString("VAR str = \"\"\nstr += ")
	} else {
		b.WriteString("RETURN ")
	}

	for i := 1; i <= 12; i++ {
		fmt.Fprintf(&b, "\"p%d-\" + %d + ", i, i)
	}

	b.WriteString("@x")

	if appendStyle {
		b.WriteString("\nRETURN str")
	}

	return b.String()
}
