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
