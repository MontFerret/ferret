package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestUdfUnusedEliminationO1(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
FUNC used() => 1
FUNC unused() => 2
RETURN used()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.UserDefined) != 1 {
				return fmt.Errorf("expected 1 UDF at O1, got %d", len(prog.Functions.UserDefined))
			}

			if prog.Functions.UserDefined[0].Name != "used" {
				return fmt.Errorf("expected used UDF, got %q", prog.Functions.UserDefined[0].Name)
			}

			return nil
		}, "unused udf eliminated"),
	}, compiler.O1)
}

func TestUdfUnusedEliminationO0(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
FUNC used() => 1
FUNC unused() => 2
RETURN used()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.UserDefined) != 2 {
				return fmt.Errorf("expected 2 UDFs at O0, got %d", len(prog.Functions.UserDefined))
			}

			return nil
		}, "unused udf metadata kept at o0"),
	})
}

func TestUdfUnusedNestedCaptureNotLifted(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
LET base = 5
FUNC outer() (
  FUNC inner() ( RETURN base )
  RETURN 1
)
RETURN outer()
`, func(prog *bytecode.Program) error {
			if len(prog.Functions.UserDefined) != 1 {
				return fmt.Errorf("expected 1 UDF at O1, got %d", len(prog.Functions.UserDefined))
			}

			outerParams := -1
			for _, udf := range prog.Functions.UserDefined {
				if udf.Name == "outer" {
					outerParams = udf.Params
					break
				}
			}

			if outerParams == -1 {
				return fmt.Errorf("expected outer UDF metadata")
			}

			if outerParams != 0 {
				return fmt.Errorf("expected outer to have 0 params (no captures), got %d", outerParams)
			}

			return nil
		}, "unused nested capture not lifted"),
	}, compiler.O1)
}

func TestUdfRecursionReachable(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
FUNC fact(n) (
  RETURN MATCH n (
    0 => 1,
    _ => n * fact(n - 1),
  )
)
RETURN fact(5)
`, func(prog *bytecode.Program) error {
			found := false
			for _, udf := range prog.Functions.UserDefined {
				if udf.Name == "fact" {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("expected fact UDF to be reachable at O1")
			}

			return nil
		}, "recursive udf remains reachable"),
	}, compiler.O1)
}
