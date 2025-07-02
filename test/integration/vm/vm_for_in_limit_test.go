package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForLimit(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LIMIT 2
				RETURN i
		`,
			[]any{1, 2}),
		CaseArray(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT 4, 2
				RETURN i
			`, []any{5, 6}),
		CaseArray(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LET x = i
				LIMIT 2
				RETURN i*x
			`, []any{1, 4},
			"Should be able to reuse values from a source"),
		CaseArray(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LET x = "foo"
				TYPENAME(x)
				LIMIT 2
				RETURN i
		`, []any{1, 2}, "Should define variables and call functions"),
		CaseArray(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT LIMIT_VALUE()
				RETURN i
		`, []any{1, 2}, "Should be able to use function call"),
		CaseArray(`
			LET o = {
				limit: 2
			}
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT o.limit
				RETURN i
		`, []any{1, 2}, "Should be able to use object property"),
		CaseArray(`
			LET o = [1,2]

			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT o[1]
				RETURN i
		`, []any{1, 2}, "Should be able to use array element"),
	}, vm.WithFunction("LIMIT_VALUE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.NewInt(2), nil
	}))
}
