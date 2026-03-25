package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestForWhileLimit(t *testing.T) {
	RunSpecs(t, []Spec{
		Array(
			`
			FOR i WHILE UNTIL(5)
				LIMIT 2
				RETURN i
		`,
			[]any{0, 1}),
		Array(`
			FOR i WHILE UNTIL(8)
				LIMIT 4, 2
				RETURN i
			`, []any{4, 5}),
		Array(`
			FOR i WHILE UNTIL(8)
				LET x = i
				LIMIT 2
				RETURN i*x
			`, []any{0, 1},
			"Should be able to reuse values from a source"),
		Array(`
			FOR i WHILE UNTIL(8)
				LET x = "foo"
				TYPENAME(x)
				LIMIT 2
				RETURN i
		`, []any{0, 1}, "Should define variables and call functions"),
		Array(`
			FOR i WHILE UNTIL(8)
				LIMIT LIMIT_VALUE()
				RETURN i
		`, []any{0, 1}, "Should be able to use function call"),
		Array(`
			LET o = {
				limit: 2
			}
			FOR i WHILE UNTIL(8)
				LIMIT o.limit
				RETURN i
		`, []any{0, 1}, "Should be able to use object property"),
		Array(`
			LET o = [1,2]

			FOR i WHILE UNTIL(8)
				LIMIT o[1]
				RETURN i
		`, []any{0, 1}, "Should be able to use array element"),
	}, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.Var().Add("LIMIT_VALUE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(2), nil
		})

		fns.From(spec.ForWhileHelpers())
	}))
}
