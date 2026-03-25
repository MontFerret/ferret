package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestForDoWhile(t *testing.T) {
	RunSpecs(t, []Spec{
		Array(`
			FOR DO WHILE false
				RETURN 1
		`, []any{1}),
		Array(`
			FOR i DO WHILE false
				RETURN i
		`, []any{0}),
		Array(`
			VAR i = 0
			FOR DO WHILE i < 2
				i = i + 1
				RETURN i - 1
		`, []any{0, 1}),
		Array(`
		FOR i DO WHILE COUNTER() < 10
				RETURN i`, []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		Array(`
			VAR i = 0
			FOR _ DO WHILE i < 1
				i = i + 1
				RETURN i
		`, []any{1}),
		Array(`
			FOR i DO WHILE UNTIL(6)
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`, []any{
			0,
			1, 2,
			2, 4, 6,
			3, 6, 9, 12,
			4, 8, 12, 16, 20,
			5, 10, 15, 20, 25, 30,
			6, 12, 18, 24, 30, 36, 42}),
	}, vm.WithFunctionsBuilder(spec.ForWhileHelpers()))
}
