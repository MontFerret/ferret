package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForWhile(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray("FOR i WHILE false RETURN i", []any{}),
		CaseArray("FOR i WHILE UNTIL(5) RETURN i", []any{0, 1, 2, 3, 4}),
		CaseArray(`
			FOR i WHILE COUNTER() < 5
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`, []any{0, 1, 2, 2, 4, 6, 3, 6, 9, 12, 4, 8, 12, 16, 20}),
	}, vm.WithFunctions(ForWhileHelpers()))
}
