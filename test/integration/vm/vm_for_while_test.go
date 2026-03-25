package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestForWhile(t *testing.T) {
	RunSpecs(t, []Spec{
		Array("FOR i WHILE false RETURN i", []any{}),
		Array("FOR i WHILE UNTIL(5) RETURN i", []any{0, 1, 2, 3, 4}),
		Array(`
			FOR i WHILE COUNTER() < 5
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`, []any{0, 1, 2, 2, 4, 6, 3, 6, 9, 12, 4, 8, 12, 16, 20}),
	}, vm.WithFunctionsBuilder(spec.ForWhileHelpers()))
}
