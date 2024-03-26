package compiler_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestFor(t *testing.T) {
	// Should not allocate memory if NONE is a return statement
	//{
	//	`FOR i IN 0..100
	//		RETURN NONE`,
	//	[]any{},
	//	ShouldEqualJSON,
	//},
	RunUseCases(t, compiler.New(), []UseCase{
		{
			"FOR i IN 1..5 RETURN i",
			[]any{1, 2, 3, 4, 5},
			ShouldEqualJSON,
		},
		{
			`FOR i IN 1..5
				                           LET x = i
				                           PRINT(x)
											RETURN i
				`,
			[]any{1, 2, 3, 4, 5},
			ShouldEqualJSON,
		},
		{
			`FOR val, counter IN 1..5
		                    LET x = val
		                    PRINT(counter)
									LET y = counter
									RETURN [x, y]
		`,
			[]any{[]any{1, 0}, []any{2, 1}, []any{3, 2}, []any{4, 3}, []any{5, 4}},
			ShouldEqualJSON,
		},
		{
			`FOR i IN [] RETURN i
		`,
			[]any{},
			ShouldEqualJSON,
		},
		{
			`FOR i IN [1, 2, 3] RETURN i
		`,
			[]any{1, 2, 3},
			ShouldEqualJSON,
		},

		{
			`FOR i, k IN [1, 2, 3] RETURN k`,
			[]any{0, 1, 2},
			ShouldEqualJSON,
		},
		{
			`FOR i IN ['foo', 'bar', 'qaz'] RETURN i`,
			[]any{"foo", "bar", "qaz"},
			ShouldEqualJSON,
		},
		{
			`FOR i IN {a: 'bar', b: 'foo', c: 'qaz'} RETURN i`,
			[]any{"foo", "bar", "qaz"},
			ShouldHaveSameItems,
		},
		{
			`FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'} RETURN k`,
			[]any{"a", "b", "c"},
			ShouldHaveSameItems,
		},
		{
			`FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
			ShouldHaveSameItems,
		},
		{
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								RETURN {[prop]: val}`,
			[]any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
			ShouldEqualJSON,
		},
		{
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								FOR val2 IN [1, 2, 3]
									RETURN { [prop]: [val, val2] }`,
			[]any{map[string]any{"a": []int{1, 1}}, map[string]any{"a": []int{1, 2}}, map[string]any{"a": []int{1, 3}}, map[string]any{"a": []int{2, 1}}, map[string]any{"a": []int{2, 2}}, map[string]any{"a": []int{2, 3}}, map[string]any{"a": []int{3, 1}}, map[string]any{"a": []int{3, 2}}, map[string]any{"a": []int{3, 3}}},
			ShouldEqualJSON,
		},
		{
			`FOR val IN [1, 2, 3]
							RETURN (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)`,
			[]any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
			ShouldEqualJSON,
		},
		{
			`FOR val IN [1, 2, 3]
							LET sub = (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)
		
							RETURN sub`,
			[]any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
			ShouldEqualJSON,
		},
		{
			`FOR i IN [ 1, 2, 3, 4, 1, 3 ]
							RETURN DISTINCT i
		`,
			[]any{1, 2, 3, 4},
			ShouldEqualJSON,
		},
	})
}
