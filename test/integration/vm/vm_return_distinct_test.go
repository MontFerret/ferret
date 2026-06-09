package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestReturnDistinct(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array("RETURN DISTINCT [3, 1, 3, 2, 1]", []any{3, 1, 2}, "preserves first occurrence order"),
		Array(`
RETURN DISTINCT [
	{ name: "Ada", role: "admin" },
	{ role: "admin", name: "Ada" }
]
`, []any{
			map[string]any{"name": "Ada", "role": "admin"},
		}, "uses Ferret object equality"),
		Array(`
LET groups = [
	["admin", "editor"],
	["editor", "viewer"]
]
RETURN DISTINCT groups[**]
`, []any{"admin", "editor", "viewer"}, "supports flattening"),
		Array(`
LET users = [
	{ role: "admin" },
	{ role: "editor" },
	{ role: "admin" }
]
RETURN DISTINCT users[*].role
`, []any{"admin", "editor"}, "supports expansion and projection"),
		Array(`
LET users = [
	{ active: true, role: "admin" },
	{ active: false, role: "viewer" },
	{ active: true, role: "admin" },
	{ active: true, role: "editor" }
]
RETURN DISTINCT users[* FILTER .active RETURN .role]
`, []any{"admin", "editor"}, "supports inline filtering and projection"),
		Array("RETURN DISTINCT []", []any{}, "supports empty arrays"),
		Array(`
FUNC unique() (
	RETURN DISTINCT [1, 2, 1]
)
RETURN unique()
`, []any{1, 2}, "supports UDF block returns"),
		Array("RETURN DISTINCT @values", []any{1, 2, 3}, "supports runtime array values").
			Env(spec.WithParam("values", []int{1, 2, 1, 3})),
		ErrorStr("RETURN DISTINCT @value", "invalid type: expected List, but got Int", "rejects runtime scalar values").
			Env(spec.WithParam("value", 42)),
	})
}

func TestReturnDistinctIdentifierEscape(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`
LET DISTINCT = [1, 2]
RETURN (DISTINCT)
`, []any{1, 2}, "returns a parenthesized DISTINCT variable"),
		Array(`
LET DISTINCT = { values: [1, 2] }
RETURN (DISTINCT.values)
`, []any{1, 2}, "returns a parenthesized DISTINCT member expression"),
		Array(`
LET DISTINCT = [[1, 2]]
RETURN (DISTINCT[0])
`, []any{1, 2}, "returns a parenthesized DISTINCT index expression"),
		Array(`
FUNC DISTINCT() => [1, 2]
RETURN (DISTINCT())
`, []any{1, 2}, "returns a parenthesized DISTINCT call expression"),
		Array(`
FOR DISTINCT IN [[1], [2]]
	RETURN (DISTINCT)
`, []any{[]any{1}, []any{2}}, "returns a parenthesized DISTINCT loop variable"),
		Array(`
FUNC read() (
	LET DISTINCT = { values: [1, 2] }
	RETURN (DISTINCT.values)
)
RETURN read()
`, []any{1, 2}, "returns a parenthesized DISTINCT expression from a UDF block"),
	})
}
