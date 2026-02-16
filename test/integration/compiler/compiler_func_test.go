package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestFunctionCall(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(
			`
			RETURN TYPENAME(1)"
		`, BC{
				I(bytecode.OpLoadConst, 2, bytecode.NewConstant(0)), // Load constant 1
				I(bytecode.OpMove, 1, 2),                            // Argument list compilation
				I(bytecode.OpType, 3, 1),                            // Call TYPENAME function
				I(bytecode.OpReturn, 3),                             // Return the result
			}),
	})
}
