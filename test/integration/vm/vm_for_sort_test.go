package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestForSort(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s
	RETURN s
`, []any{"abc", "bar", "foo", "qaz"}, "Should sort strings"),
		CaseArray(`
LET users = [
				{
					name: "Ron",
					age: 31,
					gender: "m"
				},
				{
					name: "Angela",
					age: 29,
					gender: "f"
				},
				{
					name: "Bob",
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				SORT u.name
				RETURN u
`, []any{
			map[string]any{"name": "Angela", "age": 29, "gender": "f"},
			map[string]any{"name": "Bob", "age": 36, "gender": "m"},
			map[string]any{"name": "Ron", "age": 31, "gender": "m"},
		}, "Should sort objects by name (string)"),
		CaseArray(`
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
				SORT u.age
				RETURN u
`, []any{
			map[string]any{"active": true, "age": 29, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should sort objects by age (int)"),
		CaseArray(`
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
				SORT u.age DESC
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 36, "gender": "m"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 29, "gender": "f"},
		}, "Should execute query with DESC SORT statement"),
		CaseArray(`
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
				SORT u.age ASC
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 29, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should compile query with ASC SORT statement"),
		CaseArray(`			LET users = [
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
					age: 31,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				SORT u.age, u.gender
				RETURN u`,
			[]any{
				map[string]any{"active": true, "age": 29, "gender": "f"},
				map[string]any{"active": true, "age": 31, "gender": "f"},
				map[string]any{"active": true, "age": 31, "gender": "m"},
				map[string]any{"active": true, "age": 36, "gender": "m"},
			}, "Should compile query with SORT statement with multiple expressions"),
		CaseArray(`
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
					age: 31,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				LET x = "foo"
				TEST(x)
				SORT u.age, u.gender
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 29, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should define variables and call functions"),
		CaseArray(`
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
				FILTER u.gender == "m"
				SORT u.age
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should compile query with FILTER and SORT statements"),
		CaseArray(`
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
				SORT u.age
				FILTER u.gender == "m"
				RETURN u
		`, []any{
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}, "Should compile query with SORT and FILTER statements."),
		CaseObject(`
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
			LET sorted = (FOR u IN users
				SORT u.age
				FILTER u.gender == "m"
				RETURN u)
			
			RETURN sorted[0]
		`, map[string]any{"active": true, "age": 31, "gender": "m"}, "Should return correct value from a sorted DataSet."),
	}, vm.WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.None, nil
	}))
}
