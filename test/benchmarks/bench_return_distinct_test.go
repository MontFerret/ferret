package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

const (
	returnDistinctQuery = `
RETURN DISTINCT [
	1, 2, 3, 4, 5, 6, 7, 8,
	1, 2, 3, 4, 5, 6, 7, 8
]
`

	loopDistinctQuery = `
RETURN FOR value IN [
	1, 2, 3, 4, 5, 6, 7, 8,
	1, 2, 3, 4, 5, 6, 7, 8
]
	RETURN DISTINCT value
`

	uniqueQuery = `
RETURN UNIQUE([
	1, 2, 3, 4, 5, 6, 7, 8,
	1, 2, 3, 4, 5, 6, 7, 8
])
`

	unionDistinctQuery = `
RETURN UNION_DISTINCT(
	[1, 2, 3, 4, 5, 6, 7, 8],
	[1, 2, 3, 4, 5, 6, 7, 8]
)
`

	countDistinctQuery = `
RETURN COUNT_DISTINCT([
	1, 2, 3, 4, 5, 6, 7, 8,
	1, 2, 3, 4, 5, 6, 7, 8
])
`

	compilerReturnDistinctQuery = `
FUNC unique(values) {
	RETURN DISTINCT values
}
RETURN DISTINCT unique([1, 2, 1, 3])
`
)

func BenchmarkReturnDistinct_None(b *testing.B) {
	RunBenchmarkNone(b, returnDistinctQuery)
}

func BenchmarkReturnDistinct_Full(b *testing.B) {
	RunBenchmarkFull(b, returnDistinctQuery)
}

func BenchmarkLoopDistinct_None(b *testing.B) {
	RunBenchmarkNone(b, loopDistinctQuery)
}

func BenchmarkLoopDistinct_Full(b *testing.B) {
	RunBenchmarkFull(b, loopDistinctQuery)
}

func BenchmarkUnique_None(b *testing.B) {
	RunBenchmarkNone(b, uniqueQuery)
}

func BenchmarkUnique_Full(b *testing.B) {
	RunBenchmarkFull(b, uniqueQuery)
}

func BenchmarkUnionDistinct_None(b *testing.B) {
	RunBenchmarkNone(b, unionDistinctQuery)
}

func BenchmarkUnionDistinct_Full(b *testing.B) {
	RunBenchmarkFull(b, unionDistinctQuery)
}

func BenchmarkCountDistinct_None(b *testing.B) {
	RunBenchmarkNone(b, countDistinctQuery)
}

func BenchmarkCountDistinct_Full(b *testing.B) {
	RunBenchmarkFull(b, countDistinctQuery)
}

func BenchmarkCompilerCompileReturnDistinct_None(b *testing.B) {
	benchmarkCompileQuery(b, compilerReturnDistinctQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileReturnDistinct_Full(b *testing.B) {
	benchmarkCompileQuery(b, compilerReturnDistinctQuery, compiler.OptimizationFull)
}
