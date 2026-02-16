package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestCollect(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender
				RETURN gender
`, BC{
			I(bytecode.OpReturn, 0, 7),
		}),
	})
}
