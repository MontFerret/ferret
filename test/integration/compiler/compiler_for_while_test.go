package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForWhile(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(`
			FOR i WHILE UNTIL(5)
				RETURN i
`, BC{
			I(vm.OpReturn, 0, 7),
		}),
	})
}
