package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestForNested(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(
			`FOR val IN 1..3
							FOR prop IN ["a"]
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(
			`FOR prop IN ["a"]
							FOR val IN 1..3
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								FOR val2 IN [1, 2, 3]
									RETURN { [prop]: [val, val2] }`,
			[]any{map[string]any{"a": []int{1, 1}}, map[string]any{"a": []int{1, 2}}, map[string]any{"a": []int{1, 3}}, map[string]any{"a": []int{2, 1}}, map[string]any{"a": []int{2, 2}}, map[string]any{"a": []int{2, 3}}, map[string]any{"a": []int{3, 1}}, map[string]any{"a": []int{3, 2}}, map[string]any{"a": []int{3, 3}}},
		),
		CaseArray(
			`FOR val IN [1, 2, 3]
							RETURN (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)`,
			[]any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
		),
		CaseArray(
			`FOR val IN [1, 2, 3]
							LET sub = (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)
		
							RETURN sub`,
			[]any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
		),
		CaseArray(`
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s
	FOR n IN 0..1
		RETURN CONCAT(s, n)
`, []any{"abc0", "abc1", "bar0", "bar1", "foo0", "foo1", "qaz0", "qaz1"}),
		CaseArray(`
LET strs = ["foo", "bar", "qaz", "abc"]

FOR n IN 0..1
	FOR s IN strs
		SORT s
		RETURN CONCAT(s, n)
`, []any{"abc0", "bar0", "foo0", "qaz0", "abc1", "bar1", "foo1", "qaz1"}),
		CaseArray(`
LET strs = ["foo", "bar", "qaz", "abc"]

FOR n IN 0..1
	FOR m IN 0..1
		FOR s IN strs
			SORT s
			RETURN CONCAT(s, n, m)
`, []any{"abc00", "bar00", "foo00", "qaz00", "abc01", "bar01", "foo01", "qaz01", "abc10", "bar10", "foo10", "qaz10", "abc11", "bar11", "foo11", "qaz11"}),
		CaseArray(`
LET strs = ["foo", "bar", "qaz", "abc"]

FOR n IN 0..1
	FOR s IN strs
		SORT s
		FOR m IN 0..1
			RETURN CONCAT(s, n, m)
`, []any{"abc00", "abc01", "bar00", "bar01", "foo00", "foo01", "qaz00", "qaz01", "abc10", "abc11", "bar10", "bar11", "foo10", "foo11", "qaz10", "qaz11"}),
	})
}
