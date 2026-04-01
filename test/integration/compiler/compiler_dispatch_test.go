package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestDispatchShorthandCompiles(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Opcode(`
			"click" -> @d
			RETURN 1
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile shorthand dispatch as a statement"),
		Opcode(`
			RETURN "click" -> @d
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile shorthand dispatch as an expression"),
	})
}
