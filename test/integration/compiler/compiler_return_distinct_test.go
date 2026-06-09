package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
	"github.com/MontFerret/ferret/v2/test/spec/compile/inspect"
)

func TestReturnDistinctLowering(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Opcode("RETURN DISTINCT [1, 1]", OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDistinct},
		}, "top-level RETURN DISTINCT lowers to OpDistinct"),
		Opcode("RETURN DISTINCT @values", OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDistinct},
		}, "unknown return type defers validation to runtime"),
		ProgramCheck(`
FUNC source() (
	RETURN [1, 1]
)
FUNC unique() (
	RETURN DISTINCT source()
)
RETURN unique()
`, func(prog *bytecode.Program) error {
			if !inspect.HasOpcode(prog, bytecode.OpDistinct) {
				return fmt.Errorf("expected OpDistinct in UDF return lowering")
			}

			if inspect.HasOpcode(prog, bytecode.OpTailCall) {
				return fmt.Errorf("did not expect OpTailCall for RETURN DISTINCT")
			}

			return nil
		}, "UDF RETURN DISTINCT disables tail-call lowering"),
		ProgramCheck(`
FOR value IN [1, 1]
	RETURN DISTINCT value
`, func(prog *bytecode.Program) error {
			if inspect.HasOpcode(prog, bytecode.OpDistinct) {
				return fmt.Errorf("did not expect OpDistinct for loop RETURN DISTINCT")
			}

			for _, instruction := range prog.Bytecode {
				if instruction.Opcode == bytecode.OpDataSet && instruction.Operands[1] == bytecode.Operand(1) {
					return nil
				}
			}

			return fmt.Errorf("expected distinct OpDataSet for loop RETURN DISTINCT")
		}, "loop RETURN DISTINCT keeps existing lowering"),
	})
}

func TestReturnDistinctRejectsKnownNonArrayTypes(t *testing.T) {
	const message = "RETURN DISTINCT requires an array expression"

	RunSpecs(t, []spec.Spec{
		Failure("RETURN DISTINCT 1", E{Kind: parserd.SemanticError, Message: message}),
		Failure("RETURN DISTINCT true", E{Kind: parserd.SemanticError, Message: message}),
		Failure(`RETURN DISTINCT "value"`, E{Kind: parserd.SemanticError, Message: message}),
		Failure("RETURN DISTINCT { value: 1 }", E{Kind: parserd.SemanticError, Message: message}),
		Failure("RETURN DISTINCT(1)", E{Kind: parserd.SemanticError, Message: message}, "DISTINCT immediately after RETURN is the modifier"),
		Failure(`
FUNC invalid() (
	RETURN DISTINCT 1
)
RETURN invalid()
`, E{Kind: parserd.SemanticError, Message: message}, "UDF block return rejects known scalar"),
	})
}
