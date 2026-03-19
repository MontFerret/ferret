package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestMatchLoadPropertyConstEmission(t *testing.T) {
	src := `
LET obj = { a: 1, b: 2 }
RETURN MATCH obj (
  { a: 1, b: v } => v,
  _ => 0,
)
`
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("match_load_property_const", src))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if !programHasOpcode(prog, bytecode.OpMatchLoadPropertyConst) {
		t.Fatalf("expected bytecode to contain %s", bytecode.OpMatchLoadPropertyConst)
	}

	if programHasOpcode(prog, bytecode.OpJumpIfMissingPropertyConst) {
		t.Fatalf("expected object-pattern lowering to avoid %s", bytecode.OpJumpIfMissingPropertyConst)
	}

	if programHasOpcode(prog, bytecode.OpLoadPropertyConst) {
		t.Fatalf("expected object-pattern lowering to avoid %s", bytecode.OpLoadPropertyConst)
	}

	if len(prog.Metadata.MatchFailTargets) != len(prog.Bytecode) {
		t.Fatalf("expected match fail targets metadata to align with bytecode: got %d, want %d", len(prog.Metadata.MatchFailTargets), len(prog.Bytecode))
	}
}
