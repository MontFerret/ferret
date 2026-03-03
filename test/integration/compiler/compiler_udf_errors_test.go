package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestUdfErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
FUNC f(x, x) => x
RETURN f(1)
`, E{
				Kind:    parserd.NameError,
				Message: "Parameter 'x' is already defined",
			}, "Duplicate parameter names"),
		ErrorCase(
			`
FUNC f(x) => x
RETURN f(1, 2)
`, E{
				Kind:    parserd.NameError,
				Message: "Function 'F' expects 1 arguments, got 2",
			}, "UDF wrong arity"),
	})
}
