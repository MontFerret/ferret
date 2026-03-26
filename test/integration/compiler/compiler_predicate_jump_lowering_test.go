package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
	"github.com/MontFerret/ferret/v2/test/spec/compile/inspect"
)

func TestPredicateJumpLowering_ConstEqNeLiteralSides(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected bytecode.Opcode
	}{
		{
			name: "eq literal right",
			expr: `
LET a = 1
RETURN a == 1 ? 10 : 20
`,
			expected: bytecode.OpJumpIfNeConst,
		},
		{
			name: "eq literal left",
			expr: `
LET a = 1
RETURN 1 == a ? 10 : 20
`,
			expected: bytecode.OpJumpIfNeConst,
		},
		{
			name: "ne literal right",
			expr: `
LET a = 1
RETURN a != 1 ? 10 : 20
`,
			expected: bytecode.OpJumpIfEqConst,
		},
		{
			name: "ne literal left",
			expr: `
LET a = 1
RETURN 1 != a ? 10 : 20
`,
			expected: bytecode.OpJumpIfEqConst,
		},
	}

	specs := make([]spec.Spec, 0, len(tests))
	for _, tt := range tests {
		specs = append(specs, ProgramCheck(tt.expr, func(prog *bytecode.Program) error {
			if !inspect.HasOpcode(prog, tt.expected) {
				return fmt.Errorf("expected opcode %s in lowered predicate jump", tt.expected)
			}

			return nil
		}, tt.name))
	}

	RunSpecs(t, specs)
}
