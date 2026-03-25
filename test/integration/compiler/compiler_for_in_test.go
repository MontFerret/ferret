package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestFor(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		//		SkipByteCodeCase(`
		//FOR i IN 1..5
		//	RETURN i
		//`, BC{
		//			I(bytecode.OpReturn, 0, 7),
		//		}),
	})
}
