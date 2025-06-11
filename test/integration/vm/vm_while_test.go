package vm_test

import (
	"context"
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// TODO: Implement
func TestForWhile(t *testing.T) {
	var untilCounter int
	counter := -1
	RunUseCases(t, []UseCase{
		CaseArray("FOR i WHILE false RETURN i", []any{}),
		CaseArray("FOR i WHILE UNTIL(5) RETURN i", []any{0, 1, 2, 3, 4}),
		CaseArray(`
			FOR i WHILE COUNTER() < 5
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`, []any{0, 1, 2, 2, 4, 6, 3, 6, 9, 12, 4, 8, 12, 16, 20}),
	}, vm.WithFunctions(map[string]runtime.Function{
		"UNTIL": func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			if untilCounter < int(runtime.ToIntSafe(ctx, args[0])) {
				untilCounter++

				return runtime.True, nil
			}

			return runtime.False, nil
		},
		"COUNTER": func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			counter++
			return runtime.NewInt(counter), nil
		},
	}))
}
