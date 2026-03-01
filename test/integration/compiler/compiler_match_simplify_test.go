package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestMatchSimplify_ConstantScrutinee(t *testing.T) {
	src := `
RETURN MATCH 1 (
  1 => 10,
  2 => 20,
  _ => 30,
)
`
	c0 := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	progO0, err := c0.Compile(file.NewSource("match_simplify_o0", src))
	if err != nil {
		t.Fatalf("compile O0 failed: %v", err)
	}

	c1 := compiler.New(compiler.WithOptimizationLevel(compiler.O1))
	progO1, err := c1.Compile(file.NewSource("match_simplify_o1", src))
	if err != nil {
		t.Fatalf("compile O1 failed: %v", err)
	}

	if len(progO1.Bytecode) >= len(progO0.Bytecode) {
		t.Fatalf("expected O1 bytecode to be smaller than O0: o0=%d o1=%d", len(progO0.Bytecode), len(progO1.Bytecode))
	}

	if programHasOpcode(progO1, bytecode.OpJumpIfNeConst) {
		t.Fatalf("expected simplified match to remove JumpIfNeConst")
	}
}
