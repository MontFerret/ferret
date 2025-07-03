package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForWhileFilter(t *testing.T) {
	RunUseCases(t, []UseCase{
		ByteCodeCase(`
			FOR i WHILE UNTIL(5)
				FILTER i > 2
				RETURN i
`, BC{
			I(vm.OpReturn, 0, 7),
		}),
	})
}
