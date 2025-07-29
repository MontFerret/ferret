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
				Message: "Variable 'i' is already defined",
				//Message: "Extraneous input at end of file",
			}, "Syntax error: missing return statement"),
		ErrorCase(
			`
			LET i = NONE
			RETURN
		`, E{
				Kind: compiler.SyntaxError,
				//Message: "Unexpected 'return' keyword",
				//Hint:    "Did you mean to return a value?",
			}, "Syntax error: missing return value"),
		ErrorCase(
			`
			LET i = NONE
			LET i = 1
			RETURN i
		`, E{
				Kind:    compiler.NameError,
				Message: "Variable '' is already defined",
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
