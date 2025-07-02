package vm_test

import (
	"testing"
)

func TestForCollect(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipCaseCompilationError(`
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
			FOR i IN users
				COLLECT gender = i.gender
				RETURN {
					user: i,
					gender: gender
				}
		`, "Should not have access to initial variables"),
		SkipCaseCompilationError(`
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
			FOR i IN users
				LET x = "foo"
				COLLECT gender = i.gender
				RETURN {x, gender}
		`, "Should not have access to variables defined before COLLECT"),
		CaseArray(`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender
				RETURN gender
		`,
			[]any{},
			"Should handle empty arrays gracefully"),
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
			FOR i IN users
				COLLECT gender = i.gender
				RETURN gender
`, []any{"f", "m"}, "Should group result by a single key"),
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
			FOR i IN users
				COLLECT ageGroup = FLOOR(i.age / 5)
				RETURN { ageGroup }
`, []any{
			map[string]int{"ageGroup": 5},
			map[string]int{"ageGroup": 6},
			map[string]int{"ageGroup": 7},
			map[string]int{"ageGroup": 9},
			map[string]int{"ageGroup": 13},
		}, "Should group result by a single key expression"),
		Case(`
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
			LET grouped = (FOR i IN users
				COLLECT gender = i.gender
				RETURN gender)
			RETURN grouped[0]
`, "f", "Should return correct group key by an index"),
		CaseArray(
			`LET users = [
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
			FOR i IN users
				COLLECT gender = i.gender, age = i.age
				RETURN {age, gender}
		`, []any{
				map[string]any{"age": 25, "gender": "f"},
				map[string]any{"age": 45, "gender": "f"},
				map[string]any{"age": 31, "gender": "m"},
				map[string]any{"age": 36, "gender": "m"},
				map[string]any{"age": 69, "gender": "m"},
			}, "Should group result by multiple keys"),
		CaseArray(`
LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender INTO genders
				RETURN {
					gender,
					values: genders
				}
`, []any{
			map[string]any{
				"gender": "f",
				"values": []any{
					map[string]any{
						"i": map[string]any{
							"active":  true,
							"age":     25,
							"gender":  "f",
							"married": false,
						},
					},
					map[string]any{
						"i": map[string]any{
							"active":  true,
							"age":     45,
							"gender":  "f",
							"married": true,
						},
					},
				},
			},
			map[string]any{
				"gender": "m",
				"values": []any{
					map[string]any{
						"i": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					map[string]any{
						"i": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					map[string]any{
						"i": map[string]any{
							"active":  false,
							"age":     69,
							"gender":  "m",
							"married": true,
						},
					},
				},
			},
		}, "Should create default projection"),
		CaseArray(`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender INTO genders
				RETURN {
					gender,
					values: genders
				}
`, []any{}, "COLLECT gender = i.gender INTO genders: should return an empty array when source is empty"),
		CaseArray(
			`LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender INTO genders = { active: i.active }
				RETURN {
					gender,
					values: genders
				}
		`, []any{
				map[string]any{
					"gender": "f",
					"values": []any{
						map[string]any{"active": true},
						map[string]any{"active": true},
					},
				},
				map[string]any{
					"gender": "m",
					"values": []any{
						map[string]any{"active": true},
						map[string]any{"active": true},
						map[string]any{"active": false},
					},
				},
			}, "Should create custom projection"),
		CaseArray(
			`LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender, age = i.age INTO genders = { active: i.active }
				RETURN {
					age,
					gender,
					values: genders
				}
		`, []any{
				map[string]any{
					"age":    25,
					"gender": "f",
					"values": []any{
						map[string]any{"active": true},
					},
				},
				map[string]any{
					"age":    45,
					"gender": "f",
					"values": []any{
						map[string]any{"active": true},
					},
				},
				map[string]any{
					"age":    31,
					"gender": "m",
					"values": []any{
						map[string]any{"active": true},
					},
				},
				map[string]any{
					"age":    36,
					"gender": "m",
					"values": []any{
						map[string]any{"active": true},
					},
				},
				map[string]any{
					"age":    69,
					"gender": "m",
					"values": []any{
						map[string]any{"active": false},
					},
				},
			}, "Should create custom projection grouped by multiple keys"),
		CaseArray(`
LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				LET married = i.married
				COLLECT gender = i.gender INTO genders KEEP married
				RETURN {
					gender,
					values: genders
				}
`, []any{
			map[string]any{
				"gender": "f",
				"values": []any{
					map[string]any{"married": false},
					map[string]any{"married": true},
				},
			},
			map[string]any{
				"gender": "m",
				"values": []any{
					map[string]any{"married": true},
					map[string]any{"married": false},
					map[string]any{"married": true},
				},
			},
		}, "Should create default projection with default KEEP"),
		CaseArray(`
			LET users = []
			FOR i IN users
				LET married = i.married
				COLLECT gender = i.gender INTO genders KEEP married
				RETURN {
					gender,
					values: genders
				}
`, []any{}, "COLLECT gender = i.gender INTO genders KEEP married: Should return an empty array when source is empty"),
		CaseArray(`
LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				LET married = i.married
				LET age = i.age
				COLLECT gender = i.gender INTO values KEEP married, age
				RETURN {
					gender,
					values
				}
`, []any{
			map[string]any{
				"gender": "f",
				"values": []any{
					map[string]any{
						"married": false,
						"age":     25,
					},
					map[string]any{
						"married": true,
						"age":     45,
					},
				},
			},
			map[string]any{
				"gender": "m",
				"values": []any{
					map[string]any{
						"married": true,
						"age":     31,
					},
					map[string]any{
						"married": false,
						"age":     36,
					},
					map[string]any{
						"married": true,
						"age":     69,
					},
				},
			},
		}, "Should create default projection with default KEEP using multiple keys"),
		CaseArray(`
LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				LET married = "foo"
				COLLECT gender = i.gender INTO values KEEP married
				RETURN {
					gender,
					values
				}
`, []any{
			map[string]any{
				"gender": "f",
				"values": []any{
					map[string]any{"married": "foo"},
					map[string]any{"married": "foo"},
				},
			},
			map[string]any{
				"gender": "m",
				"values": []any{
					map[string]any{"married": "foo"},
					map[string]any{"married": "foo"},
					map[string]any{"married": "foo"},
				},
			},
		}, "Should create default projection with custom KEEP"),
		CaseArray(`
LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				LET married = "foo"
				LET age = "bar"
				COLLECT gender = i.gender INTO values KEEP married, age
				RETURN {
					gender,
					values
				}
`, []any{
			map[string]any{
				"gender": "f",
				"values": []any{
					map[string]any{
						"married": "foo",
						"age":     "bar",
					},
					map[string]any{
						"married": "foo",
						"age":     "bar",
					},
				},
			},
			map[string]any{
				"gender": "m",
				"values": []any{
					map[string]any{
						"married": "foo",
						"age":     "bar",
					},
					map[string]any{
						"married": "foo",
						"age":     "bar",
					},
					map[string]any{
						"married": "foo",
						"age":     "bar",
					},
				},
			},
		}, "Should create default projection with custom KEEP using multiple keys"),
		CaseArray(`
LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				LET bar = "foo"
				COLLECT gender = i.gender INTO values KEEP bar
				RETURN {
					gender,
					values
				}
`, []any{
			map[string]any{
				"gender": "f",
				"values": []any{
					map[string]any{"bar": "foo"},
					map[string]any{"bar": "foo"},
				},
			},
			map[string]any{
				"gender": "m",
				"values": []any{
					map[string]any{"bar": "foo"},
					map[string]any{"bar": "foo"},
					map[string]any{"bar": "foo"},
				},
			},
		}, "Should create default projection with custom KEEP with custom name"),
		CaseArray(`
LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				LET bar = "foo"
				LET baz = "bar"
				COLLECT gender = i.gender INTO values KEEP bar, baz
				RETURN {
					gender,
					values
				}
`, []any{
			map[string]any{
				"gender": "f",
				"values": []any{
					map[string]any{"bar": "foo", "baz": "bar"},
					map[string]any{"bar": "foo", "baz": "bar"},
				},
			},
			map[string]any{
				"gender": "m",
				"values": []any{
					map[string]any{"bar": "foo", "baz": "bar"},
					map[string]any{"bar": "foo", "baz": "bar"},
					map[string]any{"bar": "foo", "baz": "bar"},
				},
			},
		}, "Should create default projection with custom KEEP with multiple custom names"),
		CaseArray(
			`LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender WITH COUNT INTO numberOfUsers
				RETURN {
					gender,
					values: numberOfUsers
				}
		`, []any{
				map[string]any{
					"gender": "f",
					"values": 2,
				},
				map[string]any{
					"gender": "m",
					"values": 3,
				},
			}, "Should group and count result by a single key"),

		CaseArray(
			`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender WITH COUNT INTO numberOfUsers
				RETURN {
					gender,
					values: numberOfUsers
				}
		`, []any{}, "COLLECT gender = i.gender WITH COUNT INTO numberOfUsers: Should return empty array when source is empty"),
		CaseArray(
			`LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT WITH COUNT INTO numberOfUsers
				RETURN numberOfUsers
		`, []any{
				5,
			}, "Should just count the number of items in the source"),
		CaseArray(
			`LET users = []
			FOR i IN users
				COLLECT WITH COUNT INTO numberOfUsers
				RETURN numberOfUsers
		`, []any{
				0,
			}, "Should return 0 when there are no items in the source"),
	})
}
