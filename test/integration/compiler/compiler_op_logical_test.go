package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestLogicalOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase("RETURN 1 AND 0", BC{
			I(vm.OpLoadConst, 1, C(0)),
			I(vm.OpJumpIfFalse),
			I(vm.OpLoadConst, 1, C(1)),
			I(vm.OpReturn, 1),
		}),
		SkipByteCodeCase("RETURN 1 OR 0", BC{
			I(vm.OpLoadConst, 1, C(0)),
			I(vm.OpJumpIfFalse),
			I(vm.OpLoadConst, 1, C(1)),
			I(vm.OpReturn, 1),
		}),
		//Case("RETURN 1 AND 1", 1),
		//Case("RETURN 2 > 1 AND 1 > 0", true),
		//Case("RETURN NONE && true", nil),
		//Case("RETURN '' && true", ""),
		//Case("RETURN true && 23", 23),
		//Case("RETURN 1 OR 0", 1),
		//Case("RETURN 0 OR 1", 1),
		//Case("RETURN 2 OR 1", 2),
		//Case("RETURN 2 > 1 OR 1 > 0", true),
		//Case("RETURN 2 < 1 OR 1 > 0", true),
		//Case("RETURN 1 || 7", 1),
		//Case("RETURN 0 || 7", 7),
		//Case("RETURN NONE || 'foo'", "foo"),
		//Case("RETURN '' || 'foo'", "foo"),
		//Case(`RETURN ERROR()? || 'boo'`, "boo"),
		//Case(`RETURN !ERROR()? && TRUE`, true),
		//Case(`LET u = { valid: false } RETURN u.valid || TRUE`, true),
	})
}
