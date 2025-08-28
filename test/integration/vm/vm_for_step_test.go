package vm_test

import (
	"testing"
)

func TestForStep(t *testing.T) {
	RunUseCases(t, []UseCase{
		// Basic STEP loop
		CaseArray("FOR i = 0 WHILE i < 5 STEP i = i + 1 RETURN i", []any{0, 1, 2, 3, 4}),
		
		// STEP loop with different increment
		CaseArray("FOR i = 2 WHILE i <= 8 STEP i = i + 2 RETURN i", []any{2, 4, 6, 8}),
		
		// STEP loop with decrement  
		CaseArray("FOR i = 10 WHILE i > 0 STEP i = i - 3 RETURN i", []any{10, 7, 4, 1}),
		
		// Empty STEP loop (condition false from start)
		CaseArray("FOR i = 10 WHILE i < 5 STEP i = i + 1 RETURN i", []any{}),
		
		// STEP loop with complex expression
		CaseArray("FOR i = 1 WHILE i <= 16 STEP i = i * 2 RETURN i", []any{1, 2, 4, 8, 16}),
	})
}