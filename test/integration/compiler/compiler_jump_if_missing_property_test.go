package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestJumpIfMissingPropertyConstEmission(t *testing.T) {
	src := `
LET obj = { a: 1, b: 2 }
RETURN MATCH obj (
  { a: 1 } => 1,
  _ => 0,
)
`
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("jump_if_missing_property_const", src))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if !programHasOpcode(prog, bytecode.OpJumpIfMissingPropertyConst) {
		t.Fatalf("expected bytecode to contain %s", bytecode.OpJumpIfMissingPropertyConst)
	}
}
