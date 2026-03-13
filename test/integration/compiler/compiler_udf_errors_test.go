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
		ErrorCase(
			`
FUNC outer() (
  FUNC f() => 1
  FUNC f() => 2
  RETURN f()
)
RETURN outer()
`, E{
				Kind:    parserd.NameError,
				Message: "Function 'F' is already defined",
			}, "Duplicate UDF names in the same scope"),
		ErrorCase(
			`
FUNC f(x) => x
FUNC outer() (
  FUNC f(x, y) => x + y
  RETURN f(1)
)
RETURN outer()
`, E{
				Kind:    parserd.NameError,
				Message: "Function 'F' expects 2 arguments, got 1",
			}, "Nested shadowed UDF arity uses the nearest scope"),
	})
}
