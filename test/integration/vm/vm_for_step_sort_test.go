package vm_test

import (
	"testing"
)

func TestForStepSort(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			FOR i = 5 WHILE i > 0 STEP i = i - 1
			SORT i DESC
			RETURN i
		`, []any{5, 4, 3, 2, 1}),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]

			FOR i = 0 WHILE i < 4 STEP i = i + 1
				SORT strs[i]
				RETURN i
`, []any{3, 1, 0, 2}),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]

			FOR i = 0 WHILE i < 4 STEP i = i + 1
				SORT i DESC
				RETURN strs[i]
`, []any{"abc", "qaz", "bar", "foo"}),
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
			FOR i = 0 WHILE i < 3 STEP i = i + 1
				LET u = users[i]
				SORT u.name
				RETURN users[i]
`, []any{
			map[string]any{"name": "Angela", "age": 29, "gender": "f"},
			map[string]any{"name": "Bob", "age": 36, "gender": "m"},
			map[string]any{"name": "Ron", "age": 31, "gender": "m"},
		}),
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
			FOR i = 0 WHILE i < 3 STEP i = i + 1
				LET u = users[i]
				SORT u.age DESC
				RETURN users[i]
		`, []any{
			map[string]any{"active": true, "age": 36, "gender": "m"},
			map[string]any{"active": true, "age": 31, "gender": "m"},
			map[string]any{"active": true, "age": 29, "gender": "f"},
		}),
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
			FOR i = 0 WHILE i < 4 STEP i = i + 1
				LET u = users[i]
				SORT u.age, u.gender
				RETURN users[i]`,
			[]any{
				map[string]any{"active": true, "age": 29, "gender": "f"},
				map[string]any{"active": true, "age": 31, "gender": "f"},
				map[string]any{"active": true, "age": 31, "gender": "m"},
				map[string]any{"active": true, "age": 36, "gender": "m"},
			}),
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
			LET sorted = (FOR i = 0 WHILE i < 3 STEP i = i + 1
				SORT users[i].age
				FILTER users[i].gender == "m"
				RETURN users[i])
			
			RETURN sorted[0]
		`, map[string]any{"active": true, "age": 31, "gender": "m"}),
	})
}
