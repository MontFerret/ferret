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

func TestMatchFold_ConstantScrutinee(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
RETURN MATCH 1 {
  1 => 10,
  2 => 20,
  _ => 30,
}
`, func(prog *bytecode.Program) error {
			if inspect.HasOpcode(prog, bytecode.OpJumpIfNeConst) {
				return fmt.Errorf("expected match folding to remove JumpIfNeConst in O0")
			}

			return nil
		}, "constant scrutinee folds match dispatch"),
		ProgramCheck(`
RETURN MATCH 1s {
  "1s" => 10,
  _ => 20,
}
`, func(prog *bytecode.Program) error {
			if inspect.HasOpcode(prog, bytecode.OpJumpIfNeConst) {
				return fmt.Errorf("expected strict Duration equality to fold MATCH fallback")
			}

			return nil
		}, "cross-type Duration constant MATCH folds to fallback"),
		ProgramCheck(`
RETURN MATCH 1s {
  "2s" => 10,
  _ => 20,
}
`, func(prog *bytecode.Program) error {
			if inspect.HasOpcode(prog, bytecode.OpJumpIfNeConst) {
				return fmt.Errorf("expected unequal Duration constant MATCH to fold to its fallback")
			}

			return nil
		}, "unequal cross-type Duration constant MATCH folds"),
		ProgramCheck(`
RETURN MATCH 1s {
  "tomorrow" => 10,
  _ => 20,
}
`, func(prog *bytecode.Program) error {
			if inspect.HasOpcode(prog, bytecode.OpJumpIfNeConst) {
				return fmt.Errorf("expected invalid Duration equality to fold to the fallback")
			}

			return nil
		}, "cross-type Duration equality folds constant MATCH fallback"),
	}, compiler.OptimizationNone, compiler.OptimizationFull)
}
