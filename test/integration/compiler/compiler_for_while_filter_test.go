package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestForWhileFilter(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		//		SkipByteCodeCase(`
		//			FOR i WHILE UNTIL(5)
		//				FILTER i > 2
		//				RETURN i
		//`, BC{
		//			I(bytecode.OpReturn, 0, 7),
		//		}),
	})
}
