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
		ProgramCheck(`RETURN @obj.foo.bar ON ERROR RETURN NONE`, expectCatchTableSize(1), "Member suppression should emit a guarded region"),
		ProgramCheck("RETURN QUERY VALUE `.items` IN @doc USING css ON ERROR RETURN NONE", expectCatchTableSize(1), "QUERY suppression should emit a guarded region"),
		ProgramCheck("DISPATCH \"click\" IN @d ON ERROR RETURN NONE\nRETURN 1", expectCatchTableSize(1), "DISPATCH suppression should emit a guarded region"),
		ProgramCheck("LET ok = WAITFOR VALUE NONE TIMEOUT 1ms ON TIMEOUT RETURN NONE ON ERROR FAIL\nRETURN ok", expectCatchTableSize(0), "Explicit timeout recovery should compile"),
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
				Hint:    "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>.",
			},
			"Missing condition in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FAIL or RETURN after 'ON ERROR'",
				Hint:    "Use ON ERROR FAIL to propagate failures or ON ERROR RETURN <expr> to supply a fallback value.",
			},
			"Missing error action in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR MAYBE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FAIL or RETURN after 'ON ERROR'",
				Hint:    "Use ON ERROR FAIL to propagate failures or ON ERROR RETURN <expr> to supply a fallback value.",
			},
			"Invalid error action in recovery tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR THROW`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected FAIL or RETURN after 'ON ERROR'",
				Hint:    "Use ON ERROR FAIL to propagate failures or ON ERROR RETURN <expr> to supply a fallback value.",
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
