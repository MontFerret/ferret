package vm_test

import (
	"testing"
)

func TestForStepSyntaxExamples(t *testing.T) {
	RunUseCases(t, []UseCase{
		// Examples from the problem statement
		CaseArray("FOR i = 0 WHILE i < 10 STEP i++ RETURN i", []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
		CaseArray("FOR i = 10 WHILE i > 0 STEP i-- RETURN i", []any{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}),
	})
}