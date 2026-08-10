package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestWaitforCompilationErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(`
			LET ok = WAITFOR TRUE BACKOFF UNKNOWN
			RETURN ok
		`, E{
			Message: "Unknown BACKOFF strategy",
			Hint:    "Use one of: NONE, LINEAR, EXPONENTIAL.",
		}, "Unknown BACKOFF strategy should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE OR THROW
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "OR THROW should fail as a syntax error"),
		Failure(`
			LET ok = WAITFOR TRUE JITTER 1.5
			RETURN ok
		`, E{
			Message: "JITTER must be between 0 and 1",
			Hint:    "Use a value between 0 and 1, e.g. JITTER 0.2.",
		}, "Out-of-range JITTER should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE TIMEOUT 1e999s
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Duration literal is out of range",
			Hint:    "Use a duration value that fits within the signed nanosecond range.",
		}, "Out-of-range WAITFOR TIMEOUT duration should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY 1e999s
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Duration literal is out of range",
			Hint:    "Use a duration value that fits within the signed nanosecond range.",
		}, "Out-of-range WAITFOR EVERY duration should fail compilation"),
	})
}

func TestWaitforPredicateWhenCompiles(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
			RETURN WAITFOR VALUE { state: "ready" }
				WHEN .state == "ready"
				WHEN .state != "pending"
				TIMEOUT 5ms
				EVERY 1ms
				ON TIMEOUT RETURN NONE
		`, noCompilerError, "WAITFOR VALUE should compile with repeated WHEN and wait tails"),
		ProgramCheck(`
			RETURN WAITFOR EXISTS [1, 2, 3]
				WHEN LENGTH(.) >= 3
				WHEN .[0] == 1
				TIMEOUT 5ms
				EVERY 1ms
		`, noCompilerError, "WAITFOR EXISTS should compile with repeated WHEN and wait tails"),
		ProgramCheck(`
			RETURN WAITFOR NOT EXISTS []
				WHEN LENGTH(.) == 0
				WHEN . != NONE
				TIMEOUT 5ms
				EVERY 1ms
		`, noCompilerError, "WAITFOR NOT EXISTS should compile with repeated WHEN and wait tails"),
		ProgramCheck(`
			LET obs = []
			RETURN WAITFOR EVENT "test" IN obs
				WHEN .type == "match"
				WHEN BOOM(.)
				TRIGGER (
					LET local = 1
				)
				TIMEOUT 5ms
				ON TIMEOUT RETURN NONE
		`, expectHostFunction("BOOM", 1), "WAITFOR EVENT should compile repeated WHEN, trigger, and timeout tail"),
		ProgramCheck(`
			LET obs = []
			RETURN WAITFOR EVENT "test" IN obs
				WHEN .TRIGGER == "match"
				TIMEOUT 5ms
				ON TIMEOUT RETURN NONE
		`, noCompilerError, "WAITFOR EVENT should compile TRIGGER as an implicit-current property"),
		ProgramCheck(`
			LET obs = []
			RETURN WAITFOR EVENT "test" IN obs
				TRIGGER (
					LET local = 1
				)
				TIMEOUT 5ms
				ON TIMEOUT RETURN NONE
		`, noCompilerError, "WAITFOR EVENT should compile a trigger before timeout"),
		ProgramCheck(`
			LET obs = []
			LET button = {}
			RETURN WAITFOR EVENT "test" IN obs
				TRIGGER button <- "click"
				TIMEOUT 5ms
				ON TIMEOUT RETURN NONE
		`, noCompilerError, "WAITFOR EVENT should compile inline dispatch trigger before timeout"),
		ProgramCheck(`
			LET obs = []
			VAR clicked = false
			RETURN WAITFOR EVENT "test" IN obs
				TRIGGER clicked = true
				TIMEOUT 5ms
				ON TIMEOUT RETURN NONE
		`, noCompilerError, "WAITFOR EVENT should compile inline assignment trigger before timeout"),
		ProgramCheck(`
			LET obs = []
			RETURN WAITFOR EVENT "test" IN obs
				TRIGGER ()
				TIMEOUT 5ms
				ON TIMEOUT RETURN NONE
		`, noCompilerError, "WAITFOR EVENT should compile empty trigger block before timeout"),
		ProgramCheck(`
			LET TRIGGER = 1
			RETURN TRIGGER
		`, noCompilerError, "TRIGGER should compile as a safe reserved variable name"),
		ProgramCheck(`
			RETURN @TRIGGER
		`, noCompilerError, "TRIGGER should compile as a safe reserved param name"),
		ProgramCheck(`
			RETURN TRIGGER()
		`, expectHostFunction("TRIGGER", 0), "TRIGGER should compile as a safe reserved function name"),
		ProgramCheck(`
			LET obs = []
			FOR i IN [1, 2]
				WAITFOR EVENT "test" IN obs TIMEOUT 5ms ON TIMEOUT RETURN NONE
				RETURN i
		`, noCompilerError, "WAITFOR EVENT should compile as a FOR loop body statement"),
	})
}

func TestWaitforSynchronizationGroupsCompile(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`RETURN WAITFOR ANY { false true }`, noCompilerError, "WAITFOR ANY should compile"),
		ProgramCheck(`RETURN WAITFOR ALL { true true }`, noCompilerError, "WAITFOR ALL should compile"),
		ProgramCheck(`RETURN WAITFOR EXISTS ANY { [] [1] WHEN LENGTH(.) == 1 }`, noCompilerError, "WAITFOR EXISTS ANY should compile per-arm filters"),
		ProgramCheck(`RETURN WAITFOR NOT EXISTS ALL { [] WHEN LENGTH(.) == 0 {} }`, noCompilerError, "WAITFOR NOT EXISTS ALL should compile"),
		ProgramCheck(`RETURN WAITFOR VALUE ANY { NONE "ready" WHEN . == "ready" }`, noCompilerError, "WAITFOR VALUE ANY should compile"),
		ProgramCheck(`RETURN WAITFOR VALUE ALL { "a" WHEN . == "a" "b" WHEN . == "b" }`, noCompilerError, "WAITFOR VALUE ALL should compile repeated arms"),
		ProgramCheck(`
			LET first = @first
			LET second = @second
			RETURN WAITFOR EVENT ANY {
				"first" IN first OPTIONS { capture: true } WHEN .type == "first"
				"second" IN second WHEN .type == "second" WHEN .ok
			} TRIGGER () TIMEOUT 5ms ON TIMEOUT RETURN NONE
		`, expectOpcodes(bytecode.OpStreamGroup, bytecode.OpStreamIter), "WAITFOR EVENT ANY should compile grouped descriptors"),
		ProgramCheck(`
			LET first = @first
			LET second = @second
			RETURN WAITFOR EVENT ALL {
				"first" IN first
				"second" IN second
			} TIMEOUT 5ms ON TIMEOUT RETURN NONE
		`, expectOpcodes(bytecode.OpStreamGroup, bytecode.OpStreamGroupArmDone, bytecode.OpFail), "WAITFOR EVENT ALL should compile arm completion and exhaustion failure"),
		ProgramCheck(`LET ANY = true RETURN WAITFOR ANY`, noCompilerError, "ANY should remain a singular WAITFOR identifier"),
		ProgramCheck(`LET ALL = true RETURN WAITFOR ALL`, noCompilerError, "ALL should remain a singular WAITFOR identifier"),
	})
}

func TestWaitforValuePresenceLowering(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(
			`RETURN WAITFOR VALUE @candidate TIMEOUT 1ms EVERY 0ms`,
			expectElapsedWaitWithoutHostFunctions,
			"Timed WAITFOR should use the VM clock without synthetic host calls",
		),
		Opcode(`RETURN WAITFOR VALUE @candidate TIMEOUT 1ms`, OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpJumpIfNone},
			NotExists: []bytecode.Opcode{bytecode.OpExists},
		}, "WAITFOR VALUE should use NONE presence without EXISTS semantics"),
		Opcode(`RETURN WAITFOR EXISTS @candidate TIMEOUT 1ms`, OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpExists},
			NotExists: []bytecode.Opcode{bytecode.OpJumpIfNone},
		}, "WAITFOR EXISTS should preserve EXISTS semantics"),
	}, compiler.O0, compiler.O1)
}

func expectElapsedWaitWithoutHostFunctions(program *bytecode.Program) error {
	if len(program.Functions.Host) != 0 {
		return fmt.Errorf("expected no host functions, got %v", program.Functions.Host)
	}

	elapsedCount := 0
	for _, inst := range program.Bytecode {
		if inst.Opcode == bytecode.OpElapsed {
			elapsedCount++
		}
	}

	if elapsedCount != 2 {
		return fmt.Errorf("expected two elapsed clock reads, got %d", elapsedCount)
	}

	return nil
}

func noCompilerError(*bytecode.Program) error {
	return nil
}

func expectHostFunction(name string, argCount int) func(*bytecode.Program) error {
	return func(program *bytecode.Program) error {
		for _, fn := range program.Functions.Host {
			if fn.Name == name && fn.ArgCount == argCount {
				return nil
			}
		}

		return fmt.Errorf("expected host function %q with %d arguments in %v", name, argCount, program.Functions.Host)
	}
}

func expectOpcodes(expected ...bytecode.Opcode) func(*bytecode.Program) error {
	return func(program *bytecode.Program) error {
		seen := make(map[bytecode.Opcode]bool, len(expected))
		opcodes := make([]bytecode.Opcode, 0, len(program.Bytecode))
		for _, instruction := range program.Bytecode {
			seen[instruction.Opcode] = true
			opcodes = append(opcodes, instruction.Opcode)
		}

		for _, opcode := range expected {
			if !seen[opcode] {
				return fmt.Errorf("expected opcode %s in %v", opcode, opcodes)
			}
		}

		return nil
	}
}
