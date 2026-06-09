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
FOR value IN [
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
FUNC unique(values) (
	RETURN DISTINCT values
)
RETURN DISTINCT unique([1, 2, 1, 3])
`
)

func BenchmarkReturnDistinct_O0(b *testing.B) {
	RunBenchmarkO0(b, returnDistinctQuery)
}

func BenchmarkReturnDistinct_O1(b *testing.B) {
	RunBenchmarkO1(b, returnDistinctQuery)
}

func BenchmarkLoopDistinct_O0(b *testing.B) {
	RunBenchmarkO0(b, loopDistinctQuery)
}

func BenchmarkLoopDistinct_O1(b *testing.B) {
	RunBenchmarkO1(b, loopDistinctQuery)
}

func BenchmarkUnique_O0(b *testing.B) {
	RunBenchmarkO0(b, uniqueQuery)
}

func BenchmarkUnique_O1(b *testing.B) {
	RunBenchmarkO1(b, uniqueQuery)
}

func BenchmarkUnionDistinct_O0(b *testing.B) {
	RunBenchmarkO0(b, unionDistinctQuery)
}

func BenchmarkUnionDistinct_O1(b *testing.B) {
	RunBenchmarkO1(b, unionDistinctQuery)
}

func BenchmarkCountDistinct_O0(b *testing.B) {
	RunBenchmarkO0(b, countDistinctQuery)
}

func BenchmarkCountDistinct_O1(b *testing.B) {
	RunBenchmarkO1(b, countDistinctQuery)
}

func BenchmarkCompilerCompileReturnDistinct_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerReturnDistinctQuery, compiler.O0)
}

func BenchmarkCompilerCompileReturnDistinct_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerReturnDistinctQuery, compiler.O1)
}
