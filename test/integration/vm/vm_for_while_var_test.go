package vm_test

import "testing"

func TestForWhileWithVarState(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			VAR i = 0
			FOR _ WHILE i < 5
				i = i + 1
				RETURN i - 1
		`, []any{0, 1, 2, 3, 4}),
		CaseArray(`
			VAR i = 0
			FOR _ WHILE i < 10
				i = i + 2
				RETURN i - 2
		`, []any{0, 2, 4, 6, 8}),
		CaseArray(`
			VAR i = 10
			FOR _ WHILE i > 0
				i = i - 3
				RETURN i + 3
		`, []any{10, 7, 4, 1}),
		CaseArray(`
			VAR i = 0
			FOR _ WHILE i < 1
				i = i + 1
				RETURN i - 1
		`, []any{0}),
		CaseArray(`
			VAR i = 5
			FOR _ WHILE i < 5
				i = i + 1
				RETURN i - 1
		`, []any{}),
		CaseArray(`
			VAR i = 0
			FOR _ WHILE i < 3
				i = i + 1
				RETURN [i - 1, i]
		`, []any{
			[]any{0, 1},
			[]any{1, 2},
			[]any{2, 3},
		}),
		CaseArray(`
			VAR i = 1
			FOR _ WHILE i < 20
				i = i * 2
				RETURN i / 2
		`, []any{1, 2, 4, 8, 16}),
		CaseArray(`
			VAR i = 0
			FOR _ WHILE i < 3
				i = i + 1
				RETURN i
		`, []any{1, 2, 3}),
		CaseArray(`
			VAR outer = 0
			FOR _ WHILE outer < 3
				outer = outer + 1

				VAR inner = 0
				FOR _ WHILE inner < 2
					inner = inner + 1
					RETURN [outer - 1, inner - 1]
		`, []any{
			[]any{0, 0},
			[]any{0, 1},
			[]any{1, 0},
			[]any{1, 1},
			[]any{2, 0},
			[]any{2, 1},
		}),
	})
}
