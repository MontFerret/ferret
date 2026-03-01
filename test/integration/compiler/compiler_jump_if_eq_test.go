package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestJumpIfEqConstEmission(t *testing.T) {
	src := `
LET a = 1
RETURN a != 1 ? 10 : 20
`
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("jump_if_eq_const", src))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if !programHasOpcode(prog, bytecode.OpJumpIfEqConst) {
		t.Fatalf("expected bytecode to contain %s", bytecode.OpJumpIfEqConst)
	}
}

func programHasOpcode(prog *bytecode.Program, op bytecode.Opcode) bool {
	if prog == nil {
		return false
	}

	for _, inst := range prog.Bytecode {
		if inst.Opcode == op {
			return true
		}
	}

	return false
}
