package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestVariables(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(
			`
			LET i = NONE RETURN i"
		`, BC{
				I(bytecode.OpLoadNone, 1),
				I(bytecode.OpReturn, 1),
			}, "Should be possible to use multi line string"),
		SkipByteCodeCase(`
			LET a = TRUE RETURN a
`, BC{
			I(bytecode.OpLoadBool, 1, 1),
			I(bytecode.OpReturn, 1),
		}),
		SkipByteCodeCase(`
			LET a = 1 RETURN a
`, BC{
			I(bytecode.OpLoadConst, 1, bytecode.NewConstant(0)),
			I(bytecode.OpReturn, 1),
		}),
	})
}
