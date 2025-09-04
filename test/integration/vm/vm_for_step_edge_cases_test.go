package vm_test

import (
	"testing"
)

func TestForStepEdgeCases(t *testing.T) {
	RunUseCases(t, []UseCase{
		// Test with different spacing
		CaseArray("FOR i = 0 WHILE i < 3 STEP i ++ RETURN i", []any{0, 1, 2}),
		CaseArray("FOR i = 3 WHILE i > 0 STEP i -- RETURN i", []any{3, 2, 1}),
		
		// Test with no spaces
		CaseArray("FOR i = 0 WHILE i < 3 STEP i++RETURN i", []any{0, 1, 2}),
		CaseArray("FOR i = 3 WHILE i > 0 STEP i--RETURN i", []any{3, 2, 1}),
		
		// Test with different variable names
		CaseArray("FOR counter = 0 WHILE counter < 2 STEP counter++ RETURN counter", []any{0, 1}),
		CaseArray("FOR num = 2 WHILE num > 0 STEP num-- RETURN num", []any{2, 1}),
		
		// Test with longer variable names
		CaseArray("FOR iterator = 0 WHILE iterator < 2 STEP iterator++ RETURN iterator", []any{0, 1}),
		CaseArray("FOR value_index = 2 WHILE value_index > 0 STEP value_index-- RETURN value_index", []any{2, 1}),
		
		// Test mixed with original syntax in same query (should both work)
		CaseArray("LET a = (FOR i = 0 WHILE i < 2 STEP i++ RETURN i) LET b = (FOR j = 0 WHILE j < 2 STEP j = j + 1 RETURN j) RETURN [a, b]", []any{[]any{0, 1}, []any{0, 1}}),
	})
}