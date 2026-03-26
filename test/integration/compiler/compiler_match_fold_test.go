package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestMatchFold_ConstantScrutinee(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
RETURN MATCH 1 (
  1 => 10,
  2 => 20,
  _ => 30,
)
`, func(prog *bytecode.Program) error {
			if programHasOpcode(prog, bytecode.OpJumpIfNeConst) {
				return fmt.Errorf("expected match folding to remove JumpIfNeConst in O0")
			}

			return nil
		}, "constant scrutinee folds match dispatch"),
	})
}
