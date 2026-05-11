package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestVariables(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		Failure(`
FUNC test() (
  LET x = 1
  RETURN x
)

LET test = 1

RETURN NONE
`, E{
			Kind:    parserd.NameError,
			Message: "Variable 'test' is already defined",
		}, "Should fail to compile because of variable name conflict between function and variable"),
		Failure(`
FUNC outer() (
  FUNC inner() => 1
  LET inner = 2
  RETURN inner
)

RETURN outer()
`, E{
			Kind:    parserd.NameError,
			Message: "Variable 'inner' is already defined",
		}, "Should fail to compile because of nested variable name conflict between function and variable"),
	}, compiler.O0, compiler.O1)
}

func TestVariablesInnerScopeConstantShadowingCompiles(t *testing.T) {
	expr := `
LET x = 1
LET values = (
  FOR i IN [1]
    LET x = 2
    RETURN x
)
RETURN values
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		_ = compileWithLevel(t, level, expr)
	}
}

func TestVariablesStepIdentifierCompiles(t *testing.T) {
	expr := `
LET STEP = 1
RETURN STEP
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		_ = compileWithLevel(t, level, expr)
	}
}
