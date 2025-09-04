package vm_test

import (
	"testing"
)

func TestForStep(t *testing.T) {
	RunUseCases(t, []UseCase{
		// Basic STEP loop
		CaseArray("FOR i = 0 WHILE i < 5 STEP i = i + 1 RETURN i", []any{0, 1, 2, 3, 4}),

		CaseArray("FOR i = 0 WHILE i < 5 STEP i++ RETURN i", []any{0, 1, 2, 3, 4}),

		// STEP loop with 0 results
		CaseArray("FOR i = 1 WHILE i < 1 STEP i = i + 1 RETURN i", []any{}),

		// STEP loop with 1 result
		CaseArray("FOR i = 0 WHILE i < 1 STEP i = 2 RETURN i", []any{0}),

		// STEP loop with different increment
		CaseArray("FOR i = 2 WHILE i <= 8 STEP i = i + 2 RETURN i", []any{2, 4, 6, 8}),

		// STEP loop with decrement
		CaseArray("FOR i = 10 WHILE i > 0 STEP i = i - 3 RETURN i", []any{10, 7, 4, 1}),

		CaseArray("FOR i = 5 WHILE i > 0 STEP i-- RETURN i", []any{5, 4, 3, 2, 1}),

		// Empty STEP loop (condition false from start)
		CaseArray("FOR i = 10 WHILE i < 5 STEP i = i + 1 RETURN i", []any{}),

		// STEP loop with complex expression
		CaseArray("FOR i = 1 WHILE i <= 16 STEP i = i * 2 RETURN i", []any{1, 2, 4, 8, 16}),

		// STEP loop with different variable name
		CaseArray("FOR count = 100 WHILE count >= 90 STEP count = count - 5 RETURN count", []any{100, 95, 90}),

		// STEP loop with mathematical operations in body
		CaseArray(`FOR i = 1 WHILE i <= 3 STEP i = i + 1 RETURN i * 10`, []any{10, 20, 30}),

		// STEP loop with complex condition
		CaseArray("FOR i = 0 WHILE (i * i) < 25 STEP i = i + 1 RETURN i", []any{0, 1, 2, 3, 4}),

		// STEP loop with complex step expression
		CaseArray(`
			FOR x = 1 WHILE x < 100 STEP x = x * x + 1  
			RETURN x
		`, []any{1, 2, 5, 26}), // 1, 1²+1=2, 2²+1=5, 5²+1=26, 26²+1=677 (>100)

		// STEP loop counting down by different amounts
		CaseArray(`
			FOR i = 20 WHILE i >= 0 STEP i = i - (i > 10 ? 5 : 2)
			RETURN i  
		`, []any{20, 15, 10, 8, 6, 4, 2, 0}),

		// Simpler string test
		CaseArray(`
			FOR s = "a" WHILE LENGTH(s) <= 3 STEP s = CONCAT(s, "b")
			RETURN s
		`, []any{"a", "ab", "abb"}),
	})
}
