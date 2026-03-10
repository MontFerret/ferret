package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestConstantPropagation(t *testing.T) {
	RunUseCases(t, compiler.O1, []UseCase{
		OpcodeCase(`LET a = 1 + 2 RETURN a`, OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpAdd},
		}, 3, "should fold constant addition"),

		OpcodeCase(`LET a = 1 + 2 RETURN -a`, OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpAdd, bytecode.OpFlipNegative},
		}, -3, "should fold constant addition and unary minus"),

		OpcodeCase(`LET a = 10 RETURN (a - 3) * 2`, OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpSub, bytecode.OpMul},
		}, 14, "should fold chain of arithmetic operations"),

		OpcodeCase(`RETURN 1 / 0`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDiv},
		}, runtime.ErrInvalidOperation, "should not fold division by zero"),

		OpcodeCase(`RETURN 1 / "0"`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDiv},
		}, runtime.ErrInvalidOperation, "should not fold division by zero with string"),
	})
}
