package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestForTernaryWhileExpression(t *testing.T) {
	RunSpecs(t, []Spec{
		Array(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE false RETURN i*2)
		`, []any{}),
		Array(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE T::FAIL() RETURN i*2)?
		`, []any{}),
		Array(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE COUNTER() < 10 RETURN i*2)`,
			[]any{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}),
	}, vm.WithFunctionsBuilder(spec.ForWhileHelpers()))
}
