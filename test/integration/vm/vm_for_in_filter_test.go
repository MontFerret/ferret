package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForFilter(t *testing.T) {
	counterA := 0
	counterB := 0
	RunUseCases(t, []UseCase{
		CaseArray(
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 2
				RETURN i
		`,
			[]any{3, 4, 3},
		),
		CaseArray(
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 1 AND i < 4
				RETURN i
		`,
			[]any{2, 3, 3},
		),
		CaseArray(
			`
			LET users = [
				{
					age: 31,
					gender: "m",
					name: "Josh"
				},
				{
					age: 29,
					gender: "f",
					name: "Mary"
				},
				{
					age: 36,
					gender: "m",
					name: "Peter"
				}
			]
			FOR u IN users
				FILTER u.name =~ "r"
				RETURN u
		`,
			[]any{map[string]any{"age": 29, "gender": "f", "name": "Mary"}, map[string]any{"age": 36, "gender": "m", "name": "Peter"}},
		),
		CaseArray(
			`
					LET users = [
						{
							active: true,
							age: 31,
							gender: "m"
						},
						{
							active: true,
							age: 29,
							gender: "f"
						},
						{
							active: true,
							age: 36,
							gender: "m"
						}
					]
					FOR u IN users
						FILTER u.active == true
						FILTER u.age < 35
						RETURN u
				`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}, map[string]any{"active": true, "gender": "f", "age": 29}},
		),
		CaseArray(
			`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active
				RETURN u
				`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}, map[string]any{"active": true, "gender": "f", "age": 29}, map[string]any{"active": true, "gender": "m", "age": 36}},
			"Should compile query with left side expression",
		),
		CaseArray(
			`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active == true
				LIMIT 2
				FILTER u.gender == "m"
				RETURN u
		`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}},
			"Should compile query with multiple FILTER statements",
		),
		CaseArray(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER u.active AND u.married
				RETURN u
`, []any{map[string]any{"active": true, "age": 31, "gender": "m", "married": true}, map[string]any{"active": true, "age": 45, "gender": "f", "married": true}},
			"Should compile query with multiple left side expression"),
		CaseArray(`
LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER !u.active AND u.married
				RETURN u
`, []any{map[string]any{"active": false, "age": 69, "gender": "m", "married": true}},
			"Should compile query with multiple left side expression and with binary operator"),
		CaseArray(`
		LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER !u.active AND !u.married
				RETURN u
`, []any{},
			"Should compile query with multiple left side expression and with binary operator 2"),
		CaseArray(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LET x = 2
				FILTER i > x
				RETURN i + x
`, []any{5, 6, 5}),
		CaseArray(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LET x = 2
				COUNT_A()
				FILTER i > x
				COUNT_B()
				RETURN i + x
`, []any{5, 6, 5}),
	}, vm.WithFunctions(runtime.NewFunctionsFromMap(map[string]runtime.Function{
		"COUNT_A": func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			counterA++

			return runtime.None, nil
		},
		"COUNT_B": func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			counterB++

			return runtime.None, nil
		},
	})))
}
