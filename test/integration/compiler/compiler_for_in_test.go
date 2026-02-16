package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestFor(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(`
FOR i IN 1..5
	RETURN i
`, BC{
			I(bytecode.OpReturn, 0, 7),
		}),
	})
}
