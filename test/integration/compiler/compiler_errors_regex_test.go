package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestRegexErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(`
			RETURN "abc" =~ "["
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Invalid regular expression: [",
			Hint:    "Check the syntax of the regular expression.",
		}, "Invalid regex string literal should fail compilation"),
		ErrorCase(`
			RETURN "abc" =~ 1
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Invalid regular expression: 1",
			Hint:    "Check the syntax of the regular expression.",
		}, "Non-string regex literal should fail compilation"),
	})
}
