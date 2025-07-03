package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
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
