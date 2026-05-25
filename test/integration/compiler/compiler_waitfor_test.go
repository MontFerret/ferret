package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
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
			Hint:    "Use a duration value that stays within the supported range, e.g. 100ms, 2s, or 1.5m.",
		}, "Out-of-range WAITFOR TIMEOUT duration should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY 1e999s
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Duration literal is out of range",
			Hint:    "Use a duration value that stays within the supported range, e.g. 100ms, 2s, or 1.5m.",
		}, "Out-of-range WAITFOR EVERY duration should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE TIMEOUT 1e20
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Duration literal is out of range",
			Hint:    "Use a duration value that stays within the supported range, e.g. 100ms, 2s, or 1.5m.",
		}, "Out-of-range WAITFOR TIMEOUT float constant should fail compilation"),
		Failure(`
			LET ok = WAITFOR TRUE EVERY 1e20
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Duration literal is out of range",
			Hint:    "Use a duration value that stays within the supported range, e.g. 100ms, 2s, or 1.5m.",
		}, "Out-of-range WAITFOR EVERY float constant should fail compilation"),
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
				TIMEOUT 5ms
				ON TIMEOUT RETURN NONE
		`, expectHostFunction("BOOM", 1), "WAITFOR EVENT should compile repeated WHEN host calls and timeout tail"),
	})
}

func noCompilerError(*bytecode.Program) error {
	return nil
}

func expectHostFunction(name string, arity int) func(*bytecode.Program) error {
	return func(program *bytecode.Program) error {
		got, ok := program.Functions.Host[name]
		if !ok {
			return fmt.Errorf("expected host function %q in %v", name, program.Functions.Host)
		}
		if got != arity {
			return fmt.Errorf("expected host function %q arity %d, got %d", name, arity, got)
		}

		return nil
	}
}
