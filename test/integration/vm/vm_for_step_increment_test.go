package vm_test

import (
	"testing"
)

func TestForStepIncrement(t *testing.T) {
	RunUseCases(t, []UseCase{
		// Basic STEP loop with ++ syntax
		CaseArray("FOR i = 0 WHILE i < 5 STEP i++ RETURN i", []any{0, 1, 2, 3, 4}),

		// Basic STEP loop with -- syntax
		CaseArray("FOR i = 10 WHILE i > 5 STEP i-- RETURN i", []any{10, 9, 8, 7, 6}),

		// STEP loop with ++ and 0 results
		CaseArray("FOR i = 1 WHILE i < 1 STEP i++ RETURN i", []any{}),

		// STEP loop with ++ and 1 result
		CaseArray("FOR i = 0 WHILE i < 1 STEP i++ RETURN i", []any{0}),

		// STEP loop with -- from 10 to 0
		CaseArray("FOR i = 10 WHILE i > 0 STEP i-- RETURN i", []any{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}),

		// Empty STEP loop with ++ (condition false from start)
		CaseArray("FOR i = 10 WHILE i < 5 STEP i++ RETURN i", []any{}),

		// STEP loop with different variable name and ++
		CaseArray("FOR count = 0 WHILE count < 3 STEP count++ RETURN count", []any{0, 1, 2}),

		// STEP loop with different variable name and --
		CaseArray("FOR count = 5 WHILE count > 2 STEP count-- RETURN count", []any{5, 4, 3}),

		// STEP loop with ++ and mathematical operations in body
		CaseArray(`FOR i = 1 WHILE i <= 3 STEP i++ RETURN i * 10`, []any{10, 20, 30}),

		// STEP loop with -- and mathematical operations in body
		CaseArray(`FOR i = 3 WHILE i >= 1 STEP i-- RETURN i * 10`, []any{30, 20, 10}),

		// STEP loop with ++ and complex condition
		CaseArray("FOR i = 0 WHILE (i * i) < 25 STEP i++ RETURN i", []any{0, 1, 2, 3, 4}),

		// STEP loop with -- and complex condition
		CaseArray("FOR i = 10 WHILE (i * i) > 16 STEP i-- RETURN i", []any{10, 9, 8, 7, 6, 5}),
	})
}