package bytecode_test

import (
	"github.com/MontFerret/ferret/pkg/vm"
	"testing"
)

func TestSort(t *testing.T) {
	RunUseCases(t, []UseCase{
		ByteCodeCase(`
FOR s IN []
	SORT s
	RETURN s
`, BC{
			I(vm.OpReturn, 0, 7),
		}),
	})
}
