package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func assertNoMoveFromConstant(t *testing.T, prog *bytecode.Program) {
	t.Helper()

	for i, inst := range prog.Bytecode {
		if inst.Opcode != bytecode.OpMove {
			continue
		}

		if inst.Operands[1].IsConstant() {
			t.Fatalf("unexpected MOVE from constant at pc %d: %v", i, inst.Operands)
		}
	}
}

func TestCompilerNeverEmitsMoveFromConstantOperand(t *testing.T) {
	cases := []struct {
		name string
		expr string
	}{
		{
			name: "TopLevelReturn",
			expr: `RETURN 1`,
		},
		{
			name: "UdfArrowReturn",
			expr: `
FUNC f() => 1
RETURN f()
`,
		},
		{
			name: "UdfBlockReturn",
			expr: `
FUNC f() (
  RETURN 1
)
RETURN f()
`,
		},
		{
			name: "VariableReturn",
			expr: `
LET x = 1
RETURN x
`,
		},
		{
			name: "ImplicitCurrentReturn",
			expr: `
LET items = [{ a: 1 }]
RETURN items[* RETURN .a]
`,
		},
	}

	levels := []compiler.OptimizationLevel{compiler.O0, compiler.O1}

	for _, tc := range cases {
		for _, level := range levels {
			t.Run(fmt.Sprintf("%s_O%d", tc.name, int(level)), func(t *testing.T) {
				prog := compileWithLevel(t, level, tc.expr)
				assertNoMoveFromConstant(t, prog)
			})
		}
	}
}
