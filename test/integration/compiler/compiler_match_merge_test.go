package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestMatchMerge_PureLiteralResults(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET x = @x
RETURN MATCH x (
  1 => "same",
  2 => "same",
  _ => "other",
)
`, func(prog *bytecode.Program) error {
			if got := countLoadConstValue(prog, runtime.NewString("same")); got != 1 {
				return fmt.Errorf("expected 1 load of \"same\", got %d", got)
			}

			return nil
		}, "pure literal arms merge"),
		ProgramCheck(`
LET x = @x
RETURN MATCH x (
  1 => LENGTH([1,2]),
  2 => LENGTH([1,2]),
  _ => 0,
)
`, func(prog *bytecode.Program) error {
			if got := countOpcode(prog, bytecode.OpLength); got != 2 {
				return fmt.Errorf("expected 2 LENGTH ops, got %d", got)
			}

			return nil
		}, "impure result does not merge"),
		ProgramCheck(`
LET x = @x
RETURN MATCH x (
  1 => "same",
  2 WHEN 1 < 2 => "same",
  _ => "other",
)
`, func(prog *bytecode.Program) error {
			if got := countLoadConstValue(prog, runtime.NewString("same")); got != 2 {
				return fmt.Errorf("expected 2 loads of \"same\", got %d", got)
			}

			return nil
		}, "guarded arm does not merge"),
	})
}

func countLoadConstValue(prog *bytecode.Program, val runtime.Value) int {
	if prog == nil {
		return 0
	}

	count := 0
	for _, inst := range prog.Bytecode {
		if inst.Opcode != bytecode.OpLoadConst {
			continue
		}
		constOp := inst.Operands[1]
		if !constOp.IsConstant() {
			continue
		}
		idx := constOp.Constant()
		if idx < 0 || idx >= len(prog.Constants) {
			continue
		}
		if runtime.CompareValues(prog.Constants[idx], val) == 0 {
			count++
		}
	}

	return count
}
