package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForWhileLimit(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(
			`
			FOR i WHILE UNTIL(5)
				LIMIT 2
				RETURN i
		`,
			[]any{0, 1}),
		CaseArray(`
			FOR i WHILE UNTIL(8)
				LIMIT 4, 2
				RETURN i
			`, []any{4, 5}),
		CaseArray(`
			FOR i WHILE UNTIL(8)
				LET x = i
				LIMIT 2
				RETURN i*x
			`, []any{0, 1},
			"Should be able to reuse values from a source"),
		CaseArray(`
			FOR i WHILE UNTIL(8)
				LET x = "foo"
				TYPENAME(x)
				LIMIT 2
				RETURN i
		`, []any{0, 1}, "Should define variables and call functions"),
		CaseArray(`
			FOR i WHILE UNTIL(8)
				LIMIT LIMIT_VALUE()
				RETURN i
		`, []any{0, 1}, "Should be able to use function call"),
		CaseArray(`
			LET o = {
				limit: 2
			}
			FOR i WHILE UNTIL(8)
				LIMIT o.limit
				RETURN i
		`, []any{0, 1}, "Should be able to use object property"),
		CaseArray(`
			LET o = [1,2]

			FOR i WHILE UNTIL(8)
				LIMIT o[1]
				RETURN i
		`, []any{0, 1}, "Should be able to use array element"),
	}, vm.WithFunctionSetter(func(fns runtime.Functions) {
		fns.F().Set("LIMIT_VALUE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(2), nil
		})

		fns.SetAll(ForWhileHelpers())
	}))
}
