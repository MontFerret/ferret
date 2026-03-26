package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	spec "github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestForWhileSort(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`
			FOR i WHILE UNTIL(5)
				SORT i DESC
				RETURN i
`, []any{4, 3, 2, 1, 0}),
		Array(`
			LET strs = ["foo", "bar", "qaz", "abc"]

			FOR i WHILE UNTIL(4)
				SORT strs[i]
				RETURN i
`, []any{3, 1, 0, 2}),
		Array(`
			LET strs = ["foo", "bar", "qaz", "abc"]

			FOR i WHILE UNTIL(4)
				SORT i DESC
				RETURN strs[i]
`, []any{"abc", "qaz", "bar", "foo"}),
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
			FOR i WHILE UNTIL(3)
				LET u = users[i]
				SORT u.name
				RETURN users[i]
`, []any{
			map[string]any{"name": "Angela", "age": 29, "gender": "f"},
			map[string]any{"name": "Bob", "age": 36, "gender": "m"},
			map[string]any{"name": "Ron", "age": 31, "gender": "m"},
		}),
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
			FOR i WHILE UNTIL(3)
				LET u = users[i]
				SORT u.age DESC
				RETURN users[i]
		`, []any{
			map[string]any{"active": true, "age": 36, "gender": "m"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 29, "gender": "f"},
		}),
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
			FOR i WHILE UNTIL(4)
				LET u = users[i]
				SORT u.age, u.gender
				RETURN users[i]`,
			[]any{
				map[string]any{"active": true, "age": 29, "gender": "f"},
				map[string]any{"active": true, "age": 31, "gender": "f"},
				map[string]any{"active": true, "age": 31, "gender": "m"},
				map[string]any{"active": true, "age": 36, "gender": "m"},
			}),
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
			FOR i WHILE UNTIL(4)
				LET u = users[i]
				LET x = "foo"
				TEST(x)
				SORT u.age, u.gender
				RETURN users[i]
		`, []any{
			map[string]any{"active": true, "age": 29, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "f"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}),
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
			FOR i WHILE UNTIL(3)
				LET u = users[i]
				FILTER u.gender == "m"
				SORT u.age
				RETURN users[i]
		`, []any{
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}),
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
			FOR i WHILE UNTIL(3)
				SORT users[i].age
				FILTER users[i].gender == "m"
				RETURN users[i]
		`, []any{
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 36, "gender": "m"},
		}),
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
			LET sorted = (FOR i WHILE UNTIL(3)
				SORT users[i].age
				FILTER users[i].gender == "m"
				RETURN users[i])

			RETURN sorted[0]
		`, map[string]any{"active": true, "age": 31, "gender": "m"}),
	}, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
		fns.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		})

		fns.From(spec.ForWhileHelpers())
	}))
}
