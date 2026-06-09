package benchmarks_test

import "testing"

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
