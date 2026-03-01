package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestMatchFold_ConstantScrutinee(t *testing.T) {
	src := `
RETURN MATCH 1 (
  1 => 10,
  2 => 20,
  _ => 30,
)
`
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("match_fold", src))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if programHasOpcode(prog, bytecode.OpJumpIfNeConst) {
		t.Fatalf("expected match folding to remove JumpIfNeConst in O0")
	}
}
