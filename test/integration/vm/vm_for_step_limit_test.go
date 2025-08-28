package vm_test

import (
	"testing"
)

func TestForStepLimit(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			FOR i = 1 WHILE i <= 10 STEP i = i + 1
			LIMIT 3
			RETURN i
		`, []any{1, 2, 3}),

		CaseArray(`
			FOR i = 1 WHILE i <= 10 STEP i = i + 1
			LIMIT 2, 3
			RETURN i
		`, []any{3, 4, 5}),
	})
}
