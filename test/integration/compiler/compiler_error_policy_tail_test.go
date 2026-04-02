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

func TestCompiler_ErrorPolicyTailCompiles(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`RETURN FAIL() ON ERROR SUPPRESS`, expectCatchTableSize(1), "Function call suppression should emit a guarded region"),
		ProgramCheck(`RETURN @obj.foo.bar ON ERROR SUPPRESS`, expectCatchTableSize(1), "Member suppression should emit a guarded region"),
		ProgramCheck("RETURN QUERY VALUE `.items` IN @doc USING css ON ERROR SUPPRESS", expectCatchTableSize(1), "QUERY suppression should emit a guarded region"),
		ProgramCheck("DISPATCH \"click\" IN @d ON ERROR SUPPRESS\nRETURN 1", expectCatchTableSize(1), "DISPATCH suppression should emit a guarded region"),
		ProgramCheck("LET ok = WAITFOR TRUE ON ERROR THROW\nRETURN ok", expectCatchTableSize(0), "Explicit THROW should preserve default propagation"),
		ProgramCheck("RETURN (FAIL() + 1) ON ERROR SUPPRESS", expectCatchTableSize(1), "Grouped suppression should emit a guarded region"),
	}, compiler.O0, compiler.O1)
}

func TestSyntaxErrorsErrorPolicyTail(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`RETURN FAIL() ON`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected ERROR after 'ON' in error policy tail",
				Hint:    "Complete the tail as ON ERROR SUPPRESS or ON ERROR THROW.",
			},
			"Missing ERROR in error policy tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected SUPPRESS or THROW after 'ON ERROR'",
				Hint:    "Use ON ERROR SUPPRESS to swallow failures or ON ERROR THROW to propagate them.",
			},
			"Missing action in error policy tail",
		),
		Failure(
			`RETURN FAIL() ON ERROR MAYBE`,
			E{
				Kind:    parserd.SyntaxError,
				Message: "Expected SUPPRESS or THROW after 'ON ERROR'",
				Hint:    "Use ON ERROR SUPPRESS to swallow failures or ON ERROR THROW to propagate them.",
			},
			"Invalid action in error policy tail",
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
