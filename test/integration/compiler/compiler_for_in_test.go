package bytecode_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestFor(t *testing.T) {
	RunUseCases(t, []UseCase{
		ByteCodeCase(`
FOR i IN 1..5
	RETURN i
`, BC{
			I(vm.OpReturn, 0, 7),
		}),
	})
}
