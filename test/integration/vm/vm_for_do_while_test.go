package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/MontFerret/ferret/test/integration/base"
)

func TestForDoWhile(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			FOR i DO WHILE false
				RETURN i
		`, []any{0}),
		CaseArray(`
		FOR i DO WHILE COUNTER() < 10
				RETURN i`, []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		CaseArray(`
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
	}, vm.WithFunctions(base.ForWhileHelpers()))
}
