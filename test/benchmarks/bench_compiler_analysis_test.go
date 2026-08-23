package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const compilerAnalysisQuery = `
USE WEB::HTML AS html
LET base = @base

FUNC outer(value) {
  LET carried = base + value

  FUNC middle(extra) {
    RETURN FOR item, index IN [carried, extra]
      FILTER item > 0
      RETURN html::PARSE(TO_STRING(item + index))
  }

  RETURN middle(value)
}

RETURN outer(1)
`

func BenchmarkCompilerAnalyzeSemanticSnapshot(b *testing.B) {
	compilerInstance := mustNewCompiler(
		b,
		compiler.WithOptimizationLevel(compiler.O1),
		compiler.WithDebugInfo(),
	)
	src := source.New("analysis_benchmark", compilerAnalysisQuery)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		if _, err := compilerInstance.Analyze(src); err != nil {
			b.Fatalf("analyze failed: %v", err)
		}
	}
}
