package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestForDistinct(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(
			`FOR i IN [ 1, 2, 3, 4, 1, 3 ]
							RETURN DISTINCT i
		`,
			[]any{1, 2, 3, 4},
		),
		CaseArray(
			`FOR i IN ["foo", "bar", "qaz", "foo", "abc", "bar"]
							RETURN DISTINCT i
		`,
			[]any{"foo", "bar", "qaz", "abc"},
		),
		CaseArray(
			`FOR i IN [["foo"], ["bar"], ["qaz"], ["foo"], ["abc"], ["bar"]]
							RETURN DISTINCT i
		`,
			[]any{[]any{"foo"}, []any{"bar"}, []any{"qaz"}, []any{"abc"}},
		),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "foo", "abc", "bar"]

			FOR s IN strs
				SORT s
				RETURN DISTINCT s
`, []any{"abc", "bar", "foo", "qaz"}, "Should sort and respect DISTINCT keyword"),
		CaseArray(
			`
			FOR i IN [ 1, 1, 2, 3, 4, 1, 3 ]
				LIMIT 2
				RETURN DISTINCT i
		`,
			[]any{1}),
		CaseArray(
			`
			FOR i IN [ 1, 1, 1, 3, 4, 1, 3 ]
				LIMIT 1, 2
				RETURN DISTINCT i
		`,
			[]any{1}),
		CaseArray(
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3, 3, 4 ]
				FILTER i > 2
				RETURN DISTINCT i
		`,
			[]any{3, 4},
		),
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
				RETURN DISTINCT {gender}
		`, []any{
				map[string]any{"gender": "f"},
				map[string]any{"gender": "m"},
			}),
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
				RETURN DISTINCT genders[0]
		`, []any{
				map[string]any{"active": true},
			}),
		CaseArray(`
			LET users = [
				{
					active: true,
					age: 39,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				},
				{
					active: true,
					age: 39,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 45,
					gender: "m",
					married: true
				}
			]
			FOR u IN users
			  COLLECT genderGroup = u.gender
			   AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
			
			  RETURN DISTINCT {
				minAge,
				maxAge
			  }
`, []any{
			map[string]any{"maxAge": 45, "minAge": 39},
		}, "Should collect and aggregate values by a single key"),
		// Test DISTINCT with null values
		CaseArray(`
			LET users = [
				{
					active: true,
					age: null,
					gender: "m"
				},
				{
					active: true,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					age: null,
					gender: "m"
				},
				{
					active: false,
					age: 45,
					gender: "m"
				}
			]
			FOR u IN users
				RETURN DISTINCT u.age
		`, []any{nil, 25, 45}, "Should handle null values with DISTINCT"),

		// Test DISTINCT with nested FOR loops
		CaseArray(`
			LET departments = ["IT", "Marketing", "HR"]

			FOR dept IN departments
				FOR gender IN ["m", "f"]
					RETURN DISTINCT { department: dept, gender }
		`, []any{
			map[string]any{"department": "IT", "gender": "m"},
			map[string]any{"department": "IT", "gender": "f"},
			map[string]any{"department": "Marketing", "gender": "m"},
			map[string]any{"department": "Marketing", "gender": "f"},
			map[string]any{"department": "HR", "gender": "m"},
			map[string]any{"department": "HR", "gender": "f"},
		}, "Should handle DISTINCT with nested FOR loops"),

		// Test DISTINCT with complex objects and nested properties
		CaseArray(`
			LET users = [
				{
					name: "John",
					department: {
						name: "IT",
						location: "Building A"
					}
				},
				{
					name: "Jane",
					department: {
						name: "Marketing",
						location: "Building B"
					}
				},
				{
					name: "Bob",
					department: {
						name: "IT",
						location: "Building A"
					}
				},
				{
					name: "Alice",
					department: {
						name: "HR",
						location: "Building B"
					}
				}
			]
			FOR u IN users
				RETURN DISTINCT u.department
		`, []any{
			map[string]any{"name": "IT", "location": "Building A"},
			map[string]any{"name": "Marketing", "location": "Building B"},
			map[string]any{"name": "HR", "location": "Building B"},
		}, "Should handle DISTINCT with complex objects and nested properties"),

		// Test DISTINCT with calculated values
		CaseArray(`
			LET users = [
				{ age: 25 },
				{ age: 32 },
				{ age: 45 },
				{ age: 26 },
				{ age: 31 }
			]
			FOR u IN users
				RETURN DISTINCT FLOOR(u.age / 10) * 10
		`, []any{20, 30, 40}, "Should handle DISTINCT with calculated values"),

		// Test DISTINCT with empty arrays
		CaseArray(`
			LET emptyArray = []
			FOR i IN emptyArray
				RETURN DISTINCT i
		`, []any{}, "Should handle DISTINCT with empty arrays"),

		// Test DISTINCT with SORT BY multiple fields
		CaseArray(`
			LET users = [
				{ name: "John", age: 30, gender: "m" },
				{ name: "Jane", age: 25, gender: "f" },
				{ name: "Bob", age: 30, gender: "m" },
				{ name: "Alice", age: 35, gender: "f" },
				{ name: "Mike", age: 25, gender: "m" }
			]
			FOR u IN users
				SORT u.age DESC, u.gender
				RETURN DISTINCT u.age
		`, []any{35, 30, 25}, "Should handle DISTINCT with SORT BY multiple fields"),

		// Test DISTINCT with multiple levels of nesting
		CaseArray(`
			LET departments = ["IT", "Marketing", "HR"]
			LET genders = ["m", "f"]
			LET statuses = ["active", "inactive"]

			FOR dept IN departments
				FOR gender IN genders
					FOR status IN statuses
						RETURN DISTINCT { 
							department: dept, 
							gender: gender
						}
		`, []any{
			map[string]any{"department": "IT", "gender": "m"},
			map[string]any{"department": "IT", "gender": "f"},
			map[string]any{"department": "Marketing", "gender": "m"},
			map[string]any{"department": "Marketing", "gender": "f"},
			map[string]any{"department": "HR", "gender": "m"},
			map[string]any{"department": "HR", "gender": "f"},
		}, "Should handle DISTINCT with multiple levels of nesting"),

		// Test DISTINCT with multiple levels of nesting
		CaseArray(`
			LET departments = ["IT", "Marketing", "HR"]
			LET genders = ["m", "f"]
			LET statuses = ["active", "inactive"]

			FOR dept IN departments
				SORT dept
				FOR gender IN genders
					SORT gender
					FOR status IN statuses
						SORT status
						RETURN DISTINCT { 
							department: dept, 
							gender: gender
						}
		`, []any{
			map[string]any{"department": "HR", "gender": "f"},
			map[string]any{"department": "HR", "gender": "m"},
			map[string]any{"department": "IT", "gender": "f"},
			map[string]any{"department": "IT", "gender": "m"},
			map[string]any{"department": "Marketing", "gender": "f"},
			map[string]any{"department": "Marketing", "gender": "m"},
		}, "Should handle DISTINCT with multiple levels of nesting with SORT"),

		// Test DISTINCT with a combination of COLLECT, AGGREGATE, and DISTINCT
		CaseArray(`
			LET users = [
				{ name: "John", department: "IT", age: 30 },
				{ name: "Jane", department: "Marketing", age: 25 },
				{ name: "Bob", department: "IT", age: 40 },
				{ name: "Alice", department: "HR", age: 35 },
				{ name: "Mike", department: "Marketing", age: 45 }
			]

			FOR u IN users
				COLLECT dept = u.department
				AGGREGATE avgAge = AVERAGE(u.age)
				RETURN DISTINCT {
					department: dept,
					ageCategory: avgAge > 35 ? "Senior" : "Junior"
				}
		`, []any{
			map[string]any{"department": "HR", "ageCategory": "Junior"},
			map[string]any{"department": "IT", "ageCategory": "Junior"},
			map[string]any{"department": "Marketing", "ageCategory": "Junior"},
		}, "Should handle DISTINCT with a combination of COLLECT, AGGREGATE, and DISTINCT"),

		// Test DISTINCT with array comparison and sorting
		CaseArray(`
			LET users = [
				{ name: "John", skills: ["JavaScript", "Python"] },
				{ name: "Jane", skills: ["Java", "C++"] },
				{ name: "Bob", skills: ["JavaScript", "Python"] },
				{ name: "Alice", skills: ["Python", "JavaScript"] }
			]

			FOR u IN users
				SORT u.name
				RETURN DISTINCT SORTED(u.skills)
		`, []any{
			[]any{"JavaScript", "Python"},
			[]any{"C++", "Java"},
		}, "Should handle DISTINCT with array comparison and sorting"),
	})
}
