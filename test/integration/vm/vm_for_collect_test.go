package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

// COLLECT vs. RETURN DISTINCT
//
// In order to make a result set unique, one can either use COLLECT or RETURN DISTINCT.
//
// FOR u IN users
// RETURN DISTINCT u.age
// FOR u IN users
// COLLECT age = u.age
// RETURN age
// Behind the scenes, both variants create a CollectNode. However, they use different implementations of COLLECT that have different properties:
//
// RETURN DISTINCT maintains the order of results, but it is limited to a single value.
//
// COLLECT changes the order of results (sorted or undefined), but it supports multiple values and is more flexible than RETURN DISTINCT.
//
// Aside from COLLECTs sophisticated grouping and aggregation capabilities, it allows you to place a LIMIT operation before RETURN to potentially stop the COLLECT operation early.
func TestCollect(t *testing.T) {
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
		SkipCaseArray(`
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
		SkipCase(`
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
		SkipCaseArray(
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
		SkipCaseArray(`
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
		SkipCaseArray(`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender INTO genders
				RETURN {
					gender,
					values: genders
				}
`, []any{}, "COLLECT gender = i.gender INTO genders: should return an empty array when source is empty"),
		SkipCaseArray(
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
		SkipCaseArray(
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
		SkipCaseArray(`
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
		SkipCaseArray(`
			LET users = []
			FOR i IN users
				LET married = i.married
				COLLECT gender = i.gender INTO genders KEEP married
				RETURN {
					gender,
					values: genders
				}
`, []any{}, "COLLECT gender = i.gender INTO genders KEEP married: Should return an empty array when source is empty"),
		SkipCaseArray(`
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
		SkipCaseArray(`
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
		SkipCaseArray(`
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
		SkipCaseArray(`
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
		SkipCaseArray(`
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
		SkipCaseArray(
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

		SkipCaseArray(
			`
			LET users = []
			FOR i IN users
				COLLECT gender = i.gender WITH COUNT INTO numberOfUsers
				RETURN {
					gender,
					values: numberOfUsers
				}
		`, []any{}, "COLLECT gender = i.gender WITH COUNT INTO numberOfUsers: Should return empty array when source is empty"),
		SkipCaseArray(
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
		SkipCaseArray(
			`LET users = []
			FOR i IN users
				COLLECT WITH COUNT INTO numberOfUsers
				RETURN numberOfUsers
		`, []any{
				0,
			}, "Should return 0 when there are no items in the source"),
		SkipCaseArray(`
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
FOR u IN users
  COLLECT genderGroup = u.gender
   AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)

  RETURN {
    genderGroup,
    minAge,
    maxAge
  }
`, []any{
			map[string]any{"genderGroup": "f", "minAge": 25, "maxAge": 45},
			map[string]any{"genderGroup": "m", "minAge": 31, "maxAge": 69},
		}, "Should collect and aggregate values by a single key"),
	})
}
