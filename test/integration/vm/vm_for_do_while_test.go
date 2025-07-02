package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"

	"github.com/MontFerret/ferret/pkg/vm"
)

// TODO: Implement
func TestForDoWhile(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipCaseArray(`
			FOR i DO WHILE false
				RETURN i
		`, []any{0}),
		SkipCaseArray(`
		FOR i DO WHILE COUNTER() < 10
				RETURN i`, []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		SkipCaseArray(`
			FOR i WHILE COUNTER2() < 5
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`, []any{0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 0, 2, 4, 6, 8, 0, 3, 6, 9, 12, 0, 4, 8, 12, 16}),
	}, vm.WithFunctions(ForWhileHelpers()))
}
