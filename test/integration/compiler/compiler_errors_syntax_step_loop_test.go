package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestStepLoopSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP i = i + 1 RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "STEP is no longer supported in FOR loops",
				Hint:    "Use VAR state with 'FOR _ WHILE ...' and update the counter inside the loop body.",
			}, "Legacy STEP loop is rejected"),
		ErrorCase(
			`
			FOR i = 0 WHILE STEP i = i + 1 RETURN i
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "STEP is no longer supported in FOR loops",
				Hint:    "Use VAR state with 'FOR _ WHILE ...' and update the counter inside the loop body.",
			}, "Legacy STEP syntax with missing WHILE condition is rejected"),
		ErrorCase(
			`
			LET values = (
			  FOR count = 10 WHILE count > 0 STEP count--
			    RETURN count
			)
			RETURN values
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "STEP is no longer supported in FOR loops",
				Hint:    "Use VAR state with 'FOR _ WHILE ...' and update the counter inside the loop body.",
			}, "Legacy STEP decrement syntax is rejected"),
	})
}
