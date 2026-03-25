package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestLogicalOperators(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ByteCode("RETURN 1 AND 0", BC{
			I(bytecode.OpLoadConst, 1, C(0)),
			I(bytecode.OpJumpIfFalse),
			I(bytecode.OpLoadConst, 1, C(1)),
			I(bytecode.OpReturn, 1),
		}).Skip(),
		ByteCode("RETURN 1 OR 0", BC{
			I(bytecode.OpLoadConst, 1, C(0)),
			I(bytecode.OpJumpIfFalse),
			I(bytecode.OpLoadConst, 1, C(1)),
			I(bytecode.OpReturn, 1),
		}).Skip(),
		//Spec("RETURN 1 AND 1", 1),
		//Spec("RETURN 2 > 1 AND 1 > 0", true),
		//Spec("RETURN NONE && true", nil),
		//Spec("RETURN '' && true", ""),
		//Spec("RETURN true && 23", 23),
		//Spec("RETURN 1 OR 0", 1),
		//Spec("RETURN 0 OR 1", 1),
		//Spec("RETURN 2 OR 1", 2),
		//Spec("RETURN 2 > 1 OR 1 > 0", true),
		//Spec("RETURN 2 < 1 OR 1 > 0", true),
		//Spec("RETURN 1 || 7", 1),
		//Spec("RETURN 0 || 7", 7),
		//Spec("RETURN NONE || 'foo'", "foo"),
		//Spec("RETURN '' || 'foo'", "foo"),
		//Spec(`RETURN ERROR()? || 'boo'`, "boo"),
		//Spec(`RETURN !ERROR()? && TRUE`, true),
		//Spec(`LET u = { valid: false } RETURN u.valid || TRUE`, true),
	})
}
