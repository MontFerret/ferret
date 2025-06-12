package vm_test

import (
	"context"
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestFor(t *testing.T) {
	// Should not allocate memory if NONE is a return statement
	//{
	//	`FOR i IN 0..100
	//		RETURN NONE`,
	//	[]any{},
	//	ShouldEqualJSON,
	//},
	RunUseCases(t, []UseCase{
		SkipCaseCompilationError(`
			FOR foo IN foo
				RETURN foo
		`, "Should not compile FOR foo IN foo"),
		CaseArray("FOR i IN 1..5 RETURN i", []any{1, 2, 3, 4, 5}),
		CaseArray(
			`
		FOR i IN 1..5
			LET x = i * 2
			RETURN x
		`,
			[]any{2, 4, 6, 8, 10},
		),
		CaseArray(
			`
		FOR val, counter IN 1..5
			LET x = val
			TEST_FN(counter)
			LET y = counter
			RETURN [x, y]
				`,
			[]any{[]any{1, 0}, []any{2, 1}, []any{3, 2}, []any{4, 3}, []any{5, 4}},
		),
		CaseArray(
			`FOR i IN [] RETURN i
				`,
			[]any{},
		),
		CaseArray(
			`FOR i IN [1, 2, 3] RETURN i
				`,
			[]any{1, 2, 3},
		),
		CaseArray(
			`FOR i, k IN [1, 2, 3] RETURN k`,
			[]any{0, 1, 2},
		),
		CaseArray(
			`FOR i IN ['foo', 'bar', 'qaz'] RETURN i`,
			[]any{"foo", "bar", "qaz"},
		),
		CaseItems(
			`FOR i IN {a: 'bar', b: 'foo', c: 'qaz'} RETURN i`,
			[]any{"bar", "foo", "qaz"},
		),
		CaseArray(
			`FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'} RETURN k`,
			[]any{"a", "b", "c"},
		),
		CaseArray(
			`FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
		),
		CaseArray(
			`FOR i IN { items: [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] }.items RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
		),
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
		CaseArray(
			`FOR i IN [ 1, 2, 3, 4, 1, 3 ]
							RETURN DISTINCT i
		`,
			[]any{1, 2, 3, 4},
		),
	}, vm.WithFunction("TEST_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return nil, nil
	}))
}
