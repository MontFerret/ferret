package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestFunctionCall(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(
			`
			RETURN TYPENAME(1)"
		`, BC{
				I(vm.OpLoadConst, 2, vm.NewConstant(0)), // Load constant 1
				I(vm.OpMove, 1, 2),                      // Argument list compilation
				I(vm.OpType, 3, 1),                      // Call TYPENAME function
				I(vm.OpReturn, 3),                       // Return the result
			}),
	})
}
