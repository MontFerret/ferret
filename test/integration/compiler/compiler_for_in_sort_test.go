package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestSort(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(`
FOR s IN []
	SORT s
	RETURN s
`, BC{
			I(bytecode.OpReturn, 0, 7),
		}),
	})
}
