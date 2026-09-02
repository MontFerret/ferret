package benchmarks_test

import (
	"testing"
)

const (
	matchScrutineeQuery = `
LET x = @x
RETURN MATCH x {
	0 => 0,
	1 => 1,
	_ => 2,
}
`

	matchGuardQuery = `
LET x = @x
RETURN MATCH {
	WHEN x > 10 => x,
	WHEN x > 0 => x * 2,
	_ => 0,
}
`

	matchObjectPatternQuery = `
LET obj = @obj
RETURN MATCH obj {
	{ a: 1, b: v } => v,
	_ => 0,
}
`

	matchLoopMixQuery = `
LET vals = @vals
RETURN FOR v IN vals
	RETURN MATCH v {
		0 => 0,
		1 => 1,
		2 => 2,
		_ => 3,
	}
`

	matchConstScrutineeQuery = `
RETURN MATCH 1 {
	1 => 10,
	2 => 20,
	_ => 30,
}`

	matchMergePureLiteralResults = `
LET x = @x
RETURN MATCH x {
	0 => 0,
	1 => 1,
	_ => 2,
}
`
)

var matchLoopVals = func() []any {
	vals := make([]any, 1024)
	for i := range vals {
		vals[i] = i % 4
	}
	return vals
}()

func BenchmarkMatch_Scrutinee_None(b *testing.B) {
	RunBenchmarkNone(b, matchScrutineeQuery, WithParam("x", 1))
}

func BenchmarkMatch_Scrutinee_Full(b *testing.B) {
	RunBenchmarkFull(b, matchScrutineeQuery, WithParam("x", 1))
}

func BenchmarkMatch_Guard_None(b *testing.B) {
	RunBenchmarkNone(b, matchGuardQuery, WithParam("x", 7))
}

func BenchmarkMatch_Guard_Full(b *testing.B) {
	RunBenchmarkFull(b, matchGuardQuery, WithParam("x", 7))
}

func BenchmarkMatch_ObjectPattern_None(b *testing.B) {
	RunBenchmarkNone(b, matchObjectPatternQuery, WithParam("obj", map[string]any{"a": 1, "b": 2}))
}

func BenchmarkMatch_ObjectPattern_Full(b *testing.B) {
	RunBenchmarkFull(b, matchObjectPatternQuery, WithParam("obj", map[string]any{"a": 1, "b": 2}))
}

func BenchmarkMatch_LoopMix_None(b *testing.B) {
	RunBenchmarkNone(b, matchLoopMixQuery, WithParam("vals", matchLoopVals))
}

func BenchmarkMatch_LoopMix_Full(b *testing.B) {
	RunBenchmarkFull(b, matchLoopMixQuery, WithParam("vals", matchLoopVals))
}

func BenchmarkMatch_ConstScrutinee_None(b *testing.B) {
	RunBenchmarkNone(b, matchConstScrutineeQuery)
}

func BenchmarkMatch_ConstScrutinee_Full(b *testing.B) {
	RunBenchmarkFull(b, matchConstScrutineeQuery)
}

func BenchmarkMatch_MergePureLiteralResults_None(b *testing.B) {
	RunBenchmarkNone(b, matchMergePureLiteralResults, WithParam("x", 1))
}

func BenchmarkMatch_MergePureLiteralResults_Full(b *testing.B) {
	RunBenchmarkFull(b, matchMergePureLiteralResults, WithParam("x", 1))
}
