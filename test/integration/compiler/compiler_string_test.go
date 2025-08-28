package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestString(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(
			`
			RETURN "FOO BAR"
		`, BC{
				I(vm.OpLoadConst, 1, C(0)),
				I(vm.OpReturn, 1),
			}, "Should be possible to use multi line string"),
	})
}
