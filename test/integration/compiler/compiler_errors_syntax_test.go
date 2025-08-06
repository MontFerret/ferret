package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			LET i = NONE
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected a RETURN or FOR clause at end of query",
				Hint:    "All queries must return a value. Add a RETURN statement to complete the query.",
			}, "Missing return statement"),
		ErrorCase(
			`
			LET i = NONE
			RETURN
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after 'RETURN'",
				Hint:    "Did you forget to provide a value to return?",
			}, "Missing return value"),
		ErrorCase(
			`
			FOR i IN [1, 2, 3]
				RETURN
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after 'RETURN'",
				Hint:    "Did you forget to provide a value to return?",
			}, "Missing return value in for loop"),
		ErrorCase(
			`
			LET i =
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "Expected expression after '=' for variable 'i'",
				Hint:    "Did you forget to provide a value?",
			}, "Missing variable assignment value"),
		ErrorCase(
			`
			FOR i IN 
			RETURN i
		`, E{
				Kind:    compiler.SyntaxError,
				Message: "__FAIL__",
				Hint:    "Did you forget to provide a value?",
			}, "Missing iterable in FOR"),
	})
}
