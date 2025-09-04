package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestStepLoopSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			FOR i = 0 WHILE STEP i = i + 1 RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after 'WHILE'",
				Hint:    "STEP loops require a condition after WHILE, e.g., 'FOR i = 0 WHILE i < 10 STEP i = i + 1'.",
			}, "Missing WHILE condition in STEP loop"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected a RETURN or FOR clause at end of query",
				Hint:    "All queries must return a value. Add a RETURN statement to complete the query.",
			}, "STEP followed by RETURN (parsed as incomplete query)"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP i RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected '=' after variable in STEP clause",
				Hint:    "STEP assignments require '=', e.g., 'STEP i = i + 1'.",
			}, "Missing '=' in STEP assignment"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Incomplete STEP clause",
				Hint:    "STEP requires a complete variable assignment, e.g., 'STEP i = i + 1'.",
			}, "Incomplete STEP clause at end"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP x
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Incomplete STEP clause",
				Hint:    "STEP requires a complete variable assignment, e.g., 'STEP i = i + 1'.",
			}, "Incomplete STEP clause with variable only"),

		ErrorCase(
			`
			FOR i = 0 WHILE i < 5 STEP i = RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing expression after '=' in STEP"),
	})
}
