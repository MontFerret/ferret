package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec/assert"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestVariables(t *testing.T) {
	RunSpecs(t, []Spec{
		NewSpec(`RETURN foo`, "Should not compile if a variable not defined").Expect().CompileError(assert.ShouldNotBeNil),
		NewSpec(`
			LET foo = "bar"
			LET foo = "baz"

			RETURN foo
		`, "Should not compile if a variable is not unique").Expect().CompileError(assert.ShouldNotBeNil),
		NewSpec(`			LET _ = (FOR i IN 1..100 RETURN NONE)

			RETURN _`, "Should not allow to use ignorable variable name").Expect().CompileError(assert.ShouldNotBeNil),
		Nil(`LET i = NONE RETURN i`),
		S(`LET a = TRUE RETURN a`, true),
		S(`LET a = 1 RETURN a`, 1),
		S(`LET a = 1.1 RETURN a`, 1.1),
		S(`LET a = "foo" RETURN a`, "foo"),
		S(`LET CURRENT = 5 RETURN CURRENT`, 5, "CURRENT behaves as a normal identifier"),
		S(`LET STEP = 5 RETURN STEP`, 5, "STEP behaves as a normal identifier"),
		S(
			`
		LET a = 'foo'
		LET b = a
		RETURN a`,
			"foo",
		),
		Array(`LET i = [] RETURN i`, []any{}),
		Array(`LET i = [1, 2, 3] RETURN i`, []any{1, 2, 3}),
		Array(`LET i = [None, FALSE, "foo", 1, 1.1] RETURN i`, []any{nil, false, "foo", 1, 1.1}),
		Array(`
		LET n = None
		LET b = FALSE
		LET s = "foo"
		LET i = 1
		LET f = 1.1
		LET a = [n, b, s, i, f]
		RETURN a`, []any{nil, false, "foo", 1, 1.1}),
		Object(`LET i = {} RETURN i`, map[string]any{}),
		Object(`LET i = {a: 1, b: 2} RETURN i`, map[string]any{"a": 1, "b": 2}),
		Object(`LET i = {a: 1, b: [1]} RETURN i`, map[string]any{"a": 1, "b": []any{1}}, "Nested array in object"),
		Object(`LET i = {a: {c: 1}, b: [1]} RETURN i`,
			map[string]any{"a": map[string]any{"c": 1}, "b": []any{1}}, "Nested object in object"),
		Object(`LET i = {a: 'foo', b: 1, c: TRUE, d: [], e: {}} RETURN i`,
			map[string]any{"a": "foo", "b": 1, "c": true, "d": []any{}, "e": map[string]any{}}, "Complex object"),
		Object(`LET prop = "name" LET i = { [prop]: "foo" } RETURN i`,
			map[string]any{"name": "foo"}, "Computed property name"),
		Object(`LET name="foo" LET i = { name } RETURN i`,
			map[string]any{"name": "foo"}, "Property name shorthand"),
		Array(`LET i = [{a: {c: 1}, b: [1]}] RETURN i`,
			[]any{map[string]any{"a": map[string]any{"c": 1}, "b": []any{1}}}, "Nested object in array"),
		S("LET a = 'a' LET b = a LET c = 'c' RETURN b",
			"a", "Variable reference"),
		Array("LET i = (FOR i IN [1,2,3] RETURN i) RETURN i",
			[]any{1, 2, 3}, "arrayList comprehension"),
		Array(" LET i = { items: [1,2,3]}  FOR el IN i.items RETURN el",
			[]any{1, 2, 3}, "hashMap property access for a loop source"),
		S(`LET _ = (FOR i IN 1..100 RETURN NONE) RETURN TRUE`, true),
		Array(`
			LET src = NONE
			LET i = (FOR i IN src RETURN i)?
			RETURN i
		`,
			[]any{}, "Error handling in array comprehension"),
		Array(`
			LET x = 1
			LET values = (
			  FOR i IN [1]
			    LET x = 2
			    RETURN x
			)
			RETURN [values, x]
		`, []any{
			[]any{2},
			1,
		}, "Inner-scope constant LET shadows without leaking past scope exit"),
		S(`
			LET _ = (FOR i IN 1..100 RETURN NONE)
			LET _ = (FOR i IN 1..100 RETURN NONE)

			RETURN TRUE
		`, true),
	})
}
