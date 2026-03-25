package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/compile"
	. "github.com/MontFerret/ferret/v2/test/spec/optimize"
)

func TestWaitforFastPath(t *testing.T) {
	RunUseCases(t, compiler.O1, []spec.Spec{
		Opcode(`RETURN WAITFOR TRUE TIMEOUT 1s`, compile.OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, true, "should skip sleep for true condition"),

		Opcode(`RETURN WAITFOR FALSE TIMEOUT 10ms`, compile.OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, false, "should include sleep for false condition"),

		Opcode(`RETURN WAITFOR VALUE NONE TIMEOUT 10ms`, compile.OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, nil, "should include sleep for none value"),

		Opcode(`RETURN WAITFOR EXISTS [] TIMEOUT 10ms`, compile.OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, false, "should include sleep for empty array"),

		Opcode(`RETURN WAITFOR EXISTS [1] TIMEOUT 10ms`, compile.OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, true, "should skip sleep for non-empty array"),

		Opcode(`RETURN WAITFOR EXISTS {} TIMEOUT 10ms`, compile.OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, false, "should include sleep for empty object"),

		Opcode(`RETURN WAITFOR EXISTS { foo: 1 } TIMEOUT 10ms`, compile.OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, true, "should skip sleep for non-empty object"),

		Opcode(`RETURN WAITFOR EXISTS "" TIMEOUT 10ms`, compile.OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpSleep},
			NotExists: []bytecode.Opcode{bytecode.OpJump, bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse},
		}, false, "should include sleep for empty string"),

		Opcode(`RETURN WAITFOR EXISTS "ok" TIMEOUT 10ms`, compile.OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, true, "should skip sleep for non-empty string"),

		Opcode(`RETURN WAITFOR VALUE [1] TIMEOUT 10ms`, compile.OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, []any{float64(1)}, "should skip sleep for immediate value"),

		Opcode(`RETURN WAITFOR VALUE { foo: 1 } TIMEOUT 10ms`, compile.OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, map[string]any{"foo": float64(1)}, "should skip sleep for immediate object"),

		Opcode(`RETURN WAITFOR VALUE "ok" TIMEOUT 10ms`, compile.OpcodeExistence{
			NotExists: []bytecode.Opcode{bytecode.OpSleep},
		}, "ok", "should skip sleep for immediate string"),
	})
}
