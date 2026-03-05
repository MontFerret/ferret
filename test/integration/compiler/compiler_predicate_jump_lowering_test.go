package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog := compileWithLevel(t, compiler.O0, tt.expr)
			if !hasOpcode(prog.Bytecode, tt.expected) {
				t.Fatalf("expected opcode %s in lowered predicate jump", tt.expected)
			}
		})
	}
}
