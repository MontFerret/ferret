package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestRange(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray("RETURN 1..10", []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "Should return a range from 1 to 10"),
		CaseArray("RETURN 10..1", []any{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, "Should return a range from 10 to 1"),
		CaseArray(
			`
		LET start = 1
		LET end = 10
		RETURN start..end
		`,
			[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			"Should be able to use variables in range",
		),
		CaseArray(`
		LET start = @start
		LET end = @end
		RETURN start..end
		`,
			[]any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			"Should be able to use parameters in range",
		),
	}, vm.WithParams(map[string]runtime.Value{
		"start": runtime.NewInt(1),
		"end":   runtime.NewInt(10),
	}))
}
