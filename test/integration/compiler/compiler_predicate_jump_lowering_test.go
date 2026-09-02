package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
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
			name:     "eq literal right",
			expr:     `RETURN @value == 1 ? 10 : 20`,
			expected: bytecode.OpJumpIfNeConst,
		},
		{
			name:     "eq literal left",
			expr:     `RETURN 1 == @value ? 10 : 20`,
			expected: bytecode.OpJumpIfNeConst,
		},
		{
			name:     "ne literal right",
			expr:     `RETURN @value != 1 ? 10 : 20`,
			expected: bytecode.OpJumpIfEqConst,
		},
		{
			name:     "ne literal left",
			expr:     `RETURN 1 != @value ? 10 : 20`,
			expected: bytecode.OpJumpIfEqConst,
		},
		{
			name:     "eq dynamic operands",
			expr:     `RETURN @left == @right ? 10 : 20`,
			expected: bytecode.OpJumpIfNe,
		},
		{
			name:     "ne dynamic operands",
			expr:     `RETURN @left != @right ? 10 : 20`,
			expected: bytecode.OpJumpIfEq,
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

	RunSpecsLevels(t, specs, compiler.OptimizationNone, compiler.OptimizationFull)
}
