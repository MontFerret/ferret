package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForTernaryWhileExpression(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE false RETURN i*2)
		`, []any{}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE T::FAIL() RETURN i*2)?
		`, []any{}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i WHILE COUNTER() < 10 RETURN i*2)`,
			[]any{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}),
	}, vm.WithFunctions(ForWhileHelpers()))
}
