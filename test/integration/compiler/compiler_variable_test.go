package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestVariables(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ByteCode(
			`
			LET i = NONE RETURN i"
		`, BC{
				I(bytecode.OpLoadNone, 1),
				I(bytecode.OpReturn, 1),
			}, "Should be possible to use multi line string").Skip(),
		ByteCode(`
			LET a = TRUE RETURN a
`, BC{
			I(bytecode.OpLoadBool, 1, 1),
			I(bytecode.OpReturn, 1),
		}).Skip(),
		ByteCode(`
			LET a = 1 RETURN a
`, BC{
			I(bytecode.OpLoadConst, 1, bytecode.NewConstant(0)),
			I(bytecode.OpReturn, 1),
		}).Skip(),
	})
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
