package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestForWhileWithVarState(t *testing.T) {
	RunSpecs(t, []Spec{
		Array(`
			VAR i = 0
			FOR WHILE i < 5
				i = i + 1
				RETURN i - 1
		`, []any{0, 1, 2, 3, 4}),
		Array(`
			VAR i = 0
			FOR WHILE i < 10
				i = i + 2
				RETURN i - 2
		`, []any{0, 2, 4, 6, 8}),
		Array(`
			VAR i = 10
			FOR WHILE i > 0
				i = i - 3
				RETURN i + 3
		`, []any{10, 7, 4, 1}),
		Array(`
			VAR i = 0
			FOR WHILE i < 1
				i = i + 1
				RETURN i - 1
		`, []any{0}),
		Array(`
			VAR i = 5
			FOR WHILE i < 5
				i = i + 1
				RETURN i - 1
		`, []any{}),
		Array(`
			VAR i = 0
			FOR WHILE i < 3
				i = i + 1
				RETURN [i - 1, i]
		`, []any{
			[]any{0, 1},
			[]any{1, 2},
			[]any{2, 3},
		}),
		Array(`
			VAR i = 1
			FOR WHILE i < 20
				i = i * 2
				RETURN i / 2
		`, []any{1, 2, 4, 8, 16}),
		Array(`
			VAR i = 0
			FOR WHILE i < 3
				i = i + 1
				RETURN i
		`, []any{1, 2, 3}),
		Array(`
			VAR i = 10
			FOR n WHILE i > 0
				LET out = { n, i }
				i = i - 3
				RETURN out
		`, []any{
			map[string]any{"n": 0, "i": 10},
			map[string]any{"n": 1, "i": 7},
			map[string]any{"n": 2, "i": 4},
			map[string]any{"n": 3, "i": 1},
		}),
		Array(`
			VAR outer = 0
			FOR WHILE outer < 3
				outer = outer + 1

				VAR inner = 0
				FOR WHILE inner < 2
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
		Array(`
			VAR i = 0
			FOR _ WHILE i < 3
				i = i + 1
				RETURN i
		`, []any{1, 2, 3}),
	})
}
