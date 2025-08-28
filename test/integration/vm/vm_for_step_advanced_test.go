package vm_test

import (
	"testing"
)

func TestForStepAdvanced(t *testing.T) {
	RunUseCases(t, []UseCase{
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