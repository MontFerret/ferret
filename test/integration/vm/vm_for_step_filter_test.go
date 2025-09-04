package vm_test

import (
	"testing"
)

func TestForStepFilter(t *testing.T) {
	RunUseCases(t, []UseCase{
		// STEP loop with FILTER
		CaseArray(`
			FOR i = 0 WHILE i < 10 STEP i = i + 1
			FILTER i % 2 == 0
			RETURN i
		`, []any{0, 2, 4, 6, 8}),

		// STEP loop with SORT
		CaseArray(`
			FOR i = 5 WHILE i > 0 STEP i = i - 1
			SORT i DESC
			RETURN i
		`, []any{5, 4, 3, 2, 1}),

		// STEP loop with LIMIT
		CaseArray(`
			FOR i = 1 WHILE i <= 10 STEP i = i + 1
			LIMIT 3
			RETURN i
		`, []any{1, 2, 3}),

		// STEP loop with complex expression and filter
		CaseArray(`
			FOR i = 1 WHILE i <= 5 STEP i = i + 1
			FILTER i * i <= 9
			RETURN { num: i, square: i * i }
		`, []any{
			map[string]any{"num": 1, "square": 1},
			map[string]any{"num": 2, "square": 4},
			map[string]any{"num": 3, "square": 9},
		}),
	})
}
