package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestCollectAggregateAdditional(t *testing.T) {
	RunUseCases(t, []UseCase{
		// Test 1: Multiple aggregation functions with complex expressions

		// Test 2: Nested FOR loops with COLLECT AGGREGATE

		// Test 3: Empty array handling

		// Test 4: Null value handling
		CaseArray(`
			LET users = [
				{
					active: true,
					age: null,
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
					age: null,
					gender: "m",
					married: false
				}
			]
			FOR u IN users
				COLLECT gender = u.gender
				AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
				RETURN {
					gender,
					minAge,
					maxAge
				}
		`, []any{
			map[string]any{"gender": "f", "minAge": 25, "maxAge": 25},
			map[string]any{"gender": "m", "minAge": nil, "maxAge": nil},
		}, "Should handle null values in aggregation"),

		// Test 5: Multiple grouping keys with aggregation
		CaseArray(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true,
					department: "IT"
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false,
					department: "Marketing"
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false,
					department: "IT"
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true,
					department: "Management"
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true,
					department: "Marketing"
				}
			]
			FOR u IN users
				COLLECT 
					department = u.department,
					gender = u.gender
				AGGREGATE 
					minAge = MIN(u.age), 
					maxAge = MAX(u.age)
				RETURN {
					department,
					gender,
					minAge,
					maxAge
				}
		`, []any{
			map[string]any{"department": "IT", "gender": "m", "minAge": 31, "maxAge": 36},
			map[string]any{"department": "Management", "gender": "m", "minAge": 69, "maxAge": 69},
			map[string]any{"department": "Marketing", "gender": "f", "minAge": 25, "maxAge": 45},
		}, "Should aggregate with multiple grouping keys"),

		// Test 6: Aggregation with conditional expressions
		CaseArray(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true,
					salary: 75000
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false,
					salary: 60000
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false,
					salary: 80000
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true,
					salary: 95000
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true,
					salary: 70000
				}
			]
			FOR u IN users
				COLLECT gender = u.gender
				AGGREGATE 
					activeCount = SUM(u.active ? 1 : 0),
					marriedCount = SUM(u.married ? 1 : 0),
					highSalaryCount = SUM(u.salary > 70000 ? 1 : 0)
				RETURN {
					gender,
					activeCount,
					marriedCount,
					highSalaryCount
				}
		`, []any{
			map[string]any{
				"gender":          "f",
				"activeCount":     2,
				"marriedCount":    1,
				"highSalaryCount": 0,
			},
			map[string]any{
				"gender":          "m",
				"activeCount":     2,
				"marriedCount":    2,
				"highSalaryCount": 2,
			},
		}, "Should aggregate with conditional expressions"),

		// Test 7: Aggregation with array operations
		CaseArray(`
			LET users = [
				{
					name: "John",
					skills: ["JavaScript", "Python", "Go"]
				},
				{
					name: "Jane",
					skills: ["Java", "C++", "Python"]
				},
				{
					name: "Bob",
					skills: ["Go", "Rust"]
				},
				{
					name: "Alice",
					skills: ["JavaScript", "TypeScript"]
				}
			]
			FOR u IN users
				COLLECT AGGREGATE 
					allSkills = UNION(u.skills),
					uniqueSkillCount = COUNT_DISTINCT(u.skills)
				RETURN {
					allSkills: SORTED(allSkills),
					uniqueSkillCount
				}
		`, []any{
			map[string]any{
				"allSkills":        []any{"C++", "Go", "Java", "JavaScript", "Python", "Rust", "TypeScript"},
				"uniqueSkillCount": 7,
			},
		}, "Should aggregate with array operations"),
	})
}
