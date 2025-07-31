package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			LET i = NONE
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected a RETURN or FOR clause at end of query",
				Hint:    "All queries must return a value. Add a RETURN statement to complete the query.",
			}, "Syntax error: missing return statement"),
		ErrorCase(
			`
			LET i = NONE
			RETURN
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after 'RETURN'",
				Hint:    "Did you forget to provide a value to return?",
			}, "Syntax error: missing return value"),
		ErrorCase(
			`
			LET i =
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "_FAIL_",
				Hint:    "",
			}, "Syntax error: missing variable assignment value"),
		ErrorCase(
			`
			LET i = NONE
			LET i = 1
			RETURN i
		`, E{
				Kind:    compiler.NameError,
				Message: "Variable 'i' is already defined",
			}, "Global variable not unique"),
		ErrorCase(
			`
			RETURN i
		`, E{
				Kind:    compiler.NameError,
				Message: "Variable 'i' is not defined",
			}, "Global variable not defined"),
	})
}
