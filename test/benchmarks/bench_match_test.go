package benchmarks_test

import (
	"testing"
)

const (
	matchScrutineeQuery = `
LET x = @x
RETURN MATCH x (
	0 => 0,
	1 => 1,
	_ => 2,
)
`

	matchGuardQuery = `
LET x = @x
RETURN MATCH (
	WHEN x > 10 => x,
	WHEN x > 0 => x * 2,
	_ => 0,
)
`

	matchObjectPatternQuery = `
LET obj = @obj
RETURN MATCH obj (
	{ a: 1, b: v } => v,
	_ => 0,
)
`

	matchLoopMixQuery = `
LET vals = @vals
FOR v IN vals
	RETURN MATCH v (
		0 => 0,
		1 => 1,
		2 => 2,
		_ => 3,
	)
`

	matchConstScrutineeQuery = `
RETURN MATCH 1 (
	1 => 10,
	2 => 20,
	_ => 30,
)`

	matchMergePureLiteralResults = `
LET x = @x
RETURN MATCH x (
	0 => 0,
	1 => 1,
	_ => 2,
)
`
)

var matchLoopVals = func() []any {
	vals := make([]any, 1024)
	for i := range vals {
		vals[i] = i % 4
	}
	return vals
}()

func BenchmarkMatch_Scrutinee_O0(b *testing.B) {
	RunBenchmarkO0(b, matchScrutineeQuery, WithParam("x", 1))
}

func BenchmarkMatch_Scrutinee_O1(b *testing.B) {
	RunBenchmarkO1(b, matchScrutineeQuery, WithParam("x", 1))
}

func BenchmarkMatch_Guard_O0(b *testing.B) {
	RunBenchmarkO0(b, matchGuardQuery, WithParam("x", 7))
}

func BenchmarkMatch_Guard_O1(b *testing.B) {
	RunBenchmarkO1(b, matchGuardQuery, WithParam("x", 7))
}

func BenchmarkMatch_ObjectPattern_O0(b *testing.B) {
	RunBenchmarkO0(b, matchObjectPatternQuery, WithParam("obj", map[string]any{"a": 1, "b": 2}))
}

func BenchmarkMatch_ObjectPattern_O1(b *testing.B) {
	RunBenchmarkO1(b, matchObjectPatternQuery, WithParam("obj", map[string]any{"a": 1, "b": 2}))
}

func BenchmarkMatch_LoopMix_O0(b *testing.B) {
	RunBenchmarkO0(b, matchLoopMixQuery, WithParam("vals", matchLoopVals))
}

func BenchmarkMatch_LoopMix_O1(b *testing.B) {
	RunBenchmarkO1(b, matchLoopMixQuery, WithParam("vals", matchLoopVals))
}

func BenchmarkMatch_ConstScrutinee_O0(b *testing.B) {
	RunBenchmarkO0(b, matchConstScrutineeQuery)
}

func BenchmarkMatch_ConstScrutinee_O1(b *testing.B) {
	RunBenchmarkO1(b, matchConstScrutineeQuery)
}

func BenchmarkMatch_MergePureLiteralResults_O0(b *testing.B) {
	RunBenchmarkO0(b, matchMergePureLiteralResults, WithParam("x", 1))
}

func BenchmarkMatch_MergePureLiteralResults_O1(b *testing.B) {
	RunBenchmarkO1(b, matchMergePureLiteralResults, WithParam("x", 1))
}
