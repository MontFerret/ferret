package vm_test

import (
	"context"
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// TODO: Implement
func TestForTernaryWhileExpression(t *testing.T) {
	counter := -1
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
	}, vm.WithFunctions(runtime.NewFunctionsFromMap(map[string]runtime.Function{
		"COUNTER": func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			counter++
			return runtime.NewInt(counter), nil
		},
	})))
}
