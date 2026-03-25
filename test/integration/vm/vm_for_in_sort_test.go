package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestForSort(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`
			LET strs = ["foo", "bar", "qaz", "abc"]
			
			FOR s IN strs
				SORT s
				RETURN s
`, []any{"abc", "bar", "foo", "qaz"}, "Should sort strings"),
		Array(`
			LET strs = ["foo", "bar", "qaz", "abc"]
			FOR i IN 0..3
				LET s = strs[i]
				LEt s2 = CONCAT(s, "x")
				SORT s
				RETURN CONCAT(s, s2)
`, []any{"abcabcx", "barbarx", "foofoox", "qazqazx"}, "Should keep scope of variables in FOR loop and sort strings"),
		Array(`
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
		Array(`
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
		Array(`
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
		Array(`
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
		Array(`			LET users = [
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
		Array(`
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
		Array(`
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
		Array(`
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
		Object(`
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
