package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestString(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(
			`
			RETURN "FOO BAR"
		`, BC{
				I(bytecode.OpLoadConst, 1, C(0)),
				I(bytecode.OpReturn, 1),
			}, "Should be possible to use multi line string"),
	})
}
