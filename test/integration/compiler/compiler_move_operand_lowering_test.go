package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func noMoveFromConstant(prog *bytecode.Program) error {
	for i, inst := range prog.Bytecode {
		if inst.Opcode != bytecode.OpMove && inst.Opcode != bytecode.OpMoveTracked {
			continue
		}

		if inst.Operands[1].IsConstant() {
			return fmt.Errorf("unexpected MOVE from constant at pc %d: %v", i, inst.Operands)
		}
	}

	return nil
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

	specs := make([]spec.Spec, 0, len(cases))
	for _, tc := range cases {
		specs = append(specs, ProgramCheck(tc.expr, func(prog *bytecode.Program) error {
			return noMoveFromConstant(prog)
		}, tc.name))
	}

	RunSpecsLevels(t, specs, compiler.O0, compiler.O1)
}
