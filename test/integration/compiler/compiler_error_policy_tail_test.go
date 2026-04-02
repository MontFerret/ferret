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

func TestCompiler_RecoveryTailCompiles(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`RETURN FAIL() ON ERROR RETURN NONE`, expectCatchTableSize(1), "Function call suppression should emit a guarded region"),
		ProgramCheck(`RETURN FAIL() ON ERROR RETRY 3`, expectCatchTableSize(1), "Function call retry should emit a guarded region"),
		ProgramCheck(`RETURN FAIL() ON ERROR RETRY 3 DELAY 100ms BACKOFF EXPONENTIAL OR RETURN NONE`, expectCatchTableSize(1), "Function call retry fallback should emit a guarded region"),
		ProgramCheck(`RETURN @obj.foo.bar ON ERROR RETURN NONE`, expectCatchTableSize(1), "Member suppression should emit a guarded region"),
		ProgramCheck("RETURN QUERY VALUE `.items` IN @doc USING css ON ERROR RETURN NONE", expectCatchTableSize(1), "QUERY suppression should emit a guarded region"),
		ProgramCheck("DISPATCH \"click\" IN @d ON ERROR RETURN NONE\nRETURN 1", expectCatchTableSize(1), "DISPATCH suppression should emit a guarded region"),
		ProgramCheck("LET ok = WAITFOR VALUE NONE TIMEOUT 1ms ON TIMEOUT RETURN NONE ON ERROR FAIL\nRETURN ok", expectCatchTableSize(0), "Explicit timeout recovery should compile"),
		ProgramCheck("LET ok = WAITFOR EVENT \"test\" IN @obs TIMEOUT 1ms ON TIMEOUT RETURN NONE ON ERROR RETRY 2 DELAY 5ms OR RETURN \"error\"\nRETURN ok", expectCatchTableSize(1), "WAITFOR EVENT retry should emit a guarded region"),
		ProgramCheck("LET ok = (WAITFOR VALUE NONE TIMEOUT 1ms) ON TIMEOUT RETURN NONE\nRETURN ok", expectCatchTableSize(0), "Grouped WAITFOR timeout recovery should compile"),
		ProgramCheck("LET ok = WAITFOR TRUE ON ERROR FAIL\nRETURN ok", expectCatchTableSize(0), "Explicit FAIL should preserve default propagation"),
		ProgramCheck("RETURN (FAIL() + 1) ON ERROR RETURN NONE", expectCatchTableSize(1), "Grouped suppression should emit a guarded region"),
	}, compiler.O0, compiler.O1)
}

func TestSyntaxErrorsRecoveryTail(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`RETURN FAIL() ON`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected ERROR or TIMEOUT after 'ON' in recovery tail",
				Hint:    "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON ERROR RETRY <count>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.",
			},
			"Missing condition in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FAIL, RETURN, or RETRY after 'ON ERROR'",
				Hint:    "Use ON ERROR FAIL to propagate failures, ON ERROR RETURN <expr> to supply a fallback value, or ON ERROR RETRY <count> to retry.",
			},
			"Missing error action in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR MAYBE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FAIL, RETURN, or RETRY after 'ON ERROR'",
				Hint:    "Use ON ERROR FAIL to propagate failures, ON ERROR RETURN <expr> to supply a fallback value, or ON ERROR RETRY <count> to retry.",
			},
			"Invalid error action in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR THROW`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FAIL, RETURN, or RETRY after 'ON ERROR'",
				Hint:    "Use ON ERROR FAIL to propagate failures, ON ERROR RETURN <expr> to supply a fallback value, or ON ERROR RETRY <count> to retry.",
			},
			"Legacy THROW spelling should be rejected in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR SUPPRESS`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "SUPPRESS is not supported in recovery tails",
				Hint:    "Use ON ERROR RETURN NONE instead.",
			},
			"Legacy SUPPRESS spelling should be rejected",
		),
		Failure(
			`RETURN FAIL() ON ERROR RETURN`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after 'RETURN' in recovery tail",
				Hint:    "Provide a fallback expression, e.g. ON ERROR RETURN NONE.",
			},
			"Missing fallback expression in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR RETRY`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected retry count after 'RETRY'",
				Hint:    "Provide an integer retry count, e.g. ON ERROR RETRY 3.",
			},
			"Missing retry count in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR RETRY 3 DELAY`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected value after 'DELAY' in retry policy",
				Hint:    "Provide a duration or duration-like value, e.g. DELAY 100ms.",
			},
			"Missing retry delay value should be rejected",
		),
		Failure(
			`RETURN FAIL() ON ERROR RETRY 3 DELAY 10ms BACKOFF`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected backoff kind after 'BACKOFF' in retry policy",
				Hint:    "Use BACKOFF CONSTANT, BACKOFF LINEAR, or BACKOFF EXPONENTIAL.",
			},
			"Missing retry backoff kind should be rejected",
		),
		Failure(
			`RETURN FAIL() ON ERROR RETRY 3 DELAY 10ms BACKOFF QUADRATIC`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Unknown BACKOFF strategy",
				Hint:    "Use one of: CONSTANT, LINEAR, EXPONENTIAL.",
			},
			"Unknown retry backoff should be rejected",
		),
		Failure(
			`RETURN FAIL() ON ERROR RETRY 3 BACKOFF EXPONENTIAL`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "BACKOFF requires DELAY in retry policy",
				Hint:    "Add DELAY before BACKOFF, e.g. ON ERROR RETRY 3 DELAY 100ms BACKOFF EXPONENTIAL.",
			},
			"BACKOFF without DELAY should be rejected",
		),
		Failure(
			`RETURN FAIL() ON ERROR RETRY 3 OR MAYBE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FAIL or RETURN after 'OR' in retry fallback",
				Hint:    "Complete the retry fallback as OR FAIL or OR RETURN <expr>.",
			},
			"Invalid retry fallback action should be rejected",
		),
		Failure(
			`RETURN FAIL() ON ERROR RETRY 3 OR FAIL OR RETURN NONE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Duplicate OR fallback in retry policy",
				Hint:    "Each retry policy may define OR at most once.",
			},
			"Duplicate retry OR should be rejected",
		),
		Failure(
			`RETURN FAIL() ON ERROR FAIL OR RETURN NONE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "OR is only valid inside ON ERROR RETRY",
				Hint:    "Use OR only after ON ERROR RETRY <count>, e.g. ON ERROR RETRY 3 OR RETURN NONE.",
			},
			"OR outside retry should be rejected",
		),
		Failure(
			`RETURN DISPATCH "evt" IN target ON TIMEOUT RETURN NONE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "ON TIMEOUT is only valid on WAITFOR operations",
				Hint:    "Use ON TIMEOUT only on WAITFOR expressions that define timeout handling.",
			},
			"ON TIMEOUT should be rejected outside WAITFOR",
		),
		Failure(
			`RETURN WAITFOR VALUE ready ON TIMEOUT RETURN NONE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "ON TIMEOUT requires a TIMEOUT clause on WAITFOR",
				Hint:    "Add a TIMEOUT clause before ON TIMEOUT, e.g. WAITFOR VALUE x TIMEOUT 1s ON TIMEOUT RETURN NONE.",
			},
			"ON TIMEOUT requires explicit timeout clause",
		),
		Failure(
			`RETURN WAITFOR VALUE ready TIMEOUT 1ms ON TIMEOUT RETRY 3`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "RETRY is only valid under ON ERROR",
				Hint:    "Use ON ERROR RETRY <count> ... to retry failures, or ON TIMEOUT FAIL/RETURN for timeout handling.",
			},
			"ON TIMEOUT RETRY should be rejected",
		),
		Failure(
			`RETURN FAIL() ON ERROR FAIL ON ERROR RETURN NONE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Duplicate ON ERROR handler",
				Hint:    "Each operation may define ON ERROR at most once.",
			},
			"Duplicate ON ERROR should be rejected",
		),
		Failure(
			`RETURN WAITFOR VALUE NONE TIMEOUT 1ms ON TIMEOUT RETURN NONE ON TIMEOUT FAIL`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Duplicate ON TIMEOUT handler",
				Hint:    "Each operation may define ON TIMEOUT at most once.",
			},
			"Duplicate ON TIMEOUT should be rejected",
		),
		Failure(`RETURN maybeCall?()`, E{}, "Optional call shorthand should remain invalid in v2"),
		Failure(`RETURN items?[0]`, E{}, "Optional bracket shorthand should remain invalid in v2"),
	})
}

func expectCatchTableSize(size int) func(*bytecode.Program) error {
	return func(program *bytecode.Program) error {
		if got := len(program.CatchTable); got != size {
			return fmt.Errorf("unexpected catch table size: got %d, want %d", got, size)
		}

		return nil
	}
}
