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
			@d <- "click"
			RETURN 1
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile shorthand dispatch as a statement"),
		Opcode(`
			RETURN @d <- "click"
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile shorthand dispatch as an expression"),
		Opcode(`
			LET button = @d
			button <- "click"
			LET form = @d
			form <- "submit"
			RETURN 1
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile variable-target shorthand dispatch"),
		Opcode("LET document = @d\n"+
			"document[~ css`button`][0] <- \"click\"\n"+
			"RETURN 1", OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile member-target shorthand dispatch"),
	})
}
