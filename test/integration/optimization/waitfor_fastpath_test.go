package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func TestWaitforFastPath(t *testing.T) {
	RunUseCases(t, compiler.O1, []UseCase{
		OpcodeCase(`RETURN WAITFOR TRUE TIMEOUT 1s`, OpcodeExpectation{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, true, "should skip sleep for true condition"),

		OpcodeCase(`RETURN WAITFOR FALSE TIMEOUT 10ms`, OpcodeExpectation{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, false, "should include sleep for false condition"),

		OpcodeCase(`RETURN WAITFOR VALUE NONE TIMEOUT 10ms`, OpcodeExpectation{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, nil, "should include sleep for none value"),

		OpcodeCase(`RETURN WAITFOR EXISTS [] TIMEOUT 10ms`, OpcodeExpectation{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, false, "should include sleep for empty array"),

		OpcodeCase(`RETURN WAITFOR EXISTS [1] TIMEOUT 10ms`, OpcodeExpectation{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, true, "should skip sleep for non-empty array"),

		OpcodeCase(`RETURN WAITFOR EXISTS {} TIMEOUT 10ms`, OpcodeExpectation{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, false, "should include sleep for empty object"),

		OpcodeCase(`RETURN WAITFOR EXISTS { foo: 1 } TIMEOUT 10ms`, OpcodeExpectation{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, true, "should skip sleep for non-empty object"),

		OpcodeCase(`RETURN WAITFOR EXISTS "" TIMEOUT 10ms`, OpcodeExpectation{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, false, "should include sleep for empty string"),

		OpcodeCase(`RETURN WAITFOR EXISTS "ok" TIMEOUT 10ms`, OpcodeExpectation{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, true, "should skip sleep for non-empty string"),

		OpcodeCase(`RETURN WAITFOR VALUE [1] TIMEOUT 10ms`, OpcodeExpectation{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, []any{float64(1)}, "should skip sleep for immediate value"),

		OpcodeCase(`RETURN WAITFOR VALUE { foo: 1 } TIMEOUT 10ms`, OpcodeExpectation{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, map[string]any{"foo": float64(1)}, "should skip sleep for immediate object"),

		OpcodeCase(`RETURN WAITFOR VALUE "ok" TIMEOUT 10ms`, OpcodeExpectation{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, "ok", "should skip sleep for immediate string"),
	})
}
