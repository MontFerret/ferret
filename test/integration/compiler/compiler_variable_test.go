package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestVariables(t *testing.T) {
	RunUseCases(t, []UseCase{
		ByteCodeCase(
			`
			LET i = NONE RETURN i"
		`, BC{
				I(vm.OpLoadNone, 1),
				I(vm.OpReturn, 1),
			}, "Should be possible to use multi line string"),
		ByteCodeCase(`
			LET a = TRUE RETURN a
`, BC{
			I(vm.OpLoadBool, 1, 1),
			I(vm.OpReturn, 1),
		}),
		ByteCodeCase(`
			LET a = 1 RETURN a
`, BC{
			I(vm.OpLoadConst, 1, vm.NewConstant(0)),
			I(vm.OpReturn, 1),
		}),
	})
}
