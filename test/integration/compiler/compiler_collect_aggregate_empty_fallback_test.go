package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func TestCollectAggregateGlobalEmptyFallbackLoadsNone(t *testing.T) {
	expr := `
LET users = []
FOR u IN users
	COLLECT AGGREGATE total = COUNT(u)
	RETURN total
`

	prog := compileWithLevel(t, compiler.O0, expr)

	if !hasOpcode(prog.Bytecode, bytecode.OpJumpIfTrue) {
		t.Fatalf("expected OpJumpIfTrue for empty-aggregator branch")
	}

	if got := countOpcode(prog, bytecode.OpLoadNone); got < 1 {
		t.Fatalf("expected OpLoadNone in global aggregate empty fallback, got %d", got)
	}
}
