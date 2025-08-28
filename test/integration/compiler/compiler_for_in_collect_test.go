package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestCollect(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender
				RETURN gender
`, BC{
			I(vm.OpReturn, 0, 7),
		}),
	})
}
