package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
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
	RunBenchmarkO0(b, matchScrutineeQuery, vm.WithParam("x", 1))
}

func BenchmarkMatch_Scrutinee_O1(b *testing.B) {
	RunBenchmarkO1(b, matchScrutineeQuery, vm.WithParam("x", 1))
}

func BenchmarkMatch_Guard_O0(b *testing.B) {
	RunBenchmarkO0(b, matchGuardQuery, vm.WithParam("x", 7))
}

func BenchmarkMatch_Guard_O1(b *testing.B) {
	RunBenchmarkO1(b, matchGuardQuery, vm.WithParam("x", 7))
}

func BenchmarkMatch_ObjectPattern_O0(b *testing.B) {
	RunBenchmarkO0(b, matchObjectPatternQuery, vm.WithParam("obj", map[string]any{"a": 1, "b": 2}))
}

func BenchmarkMatch_ObjectPattern_O1(b *testing.B) {
	RunBenchmarkO1(b, matchObjectPatternQuery, vm.WithParam("obj", map[string]any{"a": 1, "b": 2}))
}

func BenchmarkMatch_LoopMix_O0(b *testing.B) {
	RunBenchmarkO0(b, matchLoopMixQuery, vm.WithParam("vals", matchLoopVals))
}

func BenchmarkMatch_LoopMix_O1(b *testing.B) {
	RunBenchmarkO1(b, matchLoopMixQuery, vm.WithParam("vals", matchLoopVals))
}

func BenchmarkMatch_ConstScrutinee_O0(b *testing.B) {
	RunBenchmarkO0(b, matchConstScrutineeQuery)
}

func BenchmarkMatch_ConstScrutinee_O1(b *testing.B) {
	RunBenchmarkO1(b, matchConstScrutineeQuery)
}
