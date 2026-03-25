package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestString(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ByteCode(
			`
			RETURN "FOO BAR"
		`, BC{
				I(bytecode.OpLoadConst, 1, C(0)),
				I(bytecode.OpReturn, 1),
			}, "Should be possible to use multi line string").Skip(),
	})
}
