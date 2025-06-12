package bytecode_test

import (
	"github.com/MontFerret/ferret/pkg/vm"
	"testing"
)

func TestCollect(t *testing.T) {
	RunUseCases(t, []UseCase{
		ByteCodeCase(`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender
				RETURN gender
`, BC{
			I(vm.OpReturn, 0, 7),
		}),
	})
}
