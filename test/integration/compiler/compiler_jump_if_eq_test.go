package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestJumpIfEqConstEmission(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET a = 1
RETURN a != 1 ? 10 : 20
`, func(prog *bytecode.Program) error {
			if !programHasOpcode(prog, bytecode.OpJumpIfEqConst) {
				return fmt.Errorf("expected bytecode to contain %s", bytecode.OpJumpIfEqConst)
			}

			return nil
		}, "ternary lowering uses JumpIfEqConst"),
	})
}
