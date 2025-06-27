package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestCollectAggregate(t *testing.T) {
	RunUseCases(t, []UseCase{
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
  				COLLECT AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  				RETURN {
    				minAge, 
    				maxAge 
  				}
		`,
			[]any{map[string]any{"minAge": 25, "maxAge": 69}},
			"Should collect and aggregate values without grouping"),
		CaseArray(`
			LET users = []
			FOR u IN users
				COLLECT AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
				RETURN {
					minAge,
					maxAge
				}
		`,
			[]any{map[string]any{"minAge": nil, "maxAge": nil}},
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
			FOR u IN users
  				COLLECT AGGREGATE ages = UNION(u.age, u.age)
  				RETURN { ages } 
`, []any{
			map[string]any{"ages": []any{31, 25, 36, 69, 45, 31, 25, 36, 69, 45}},
		}, "Should call aggregation functions with more than one argument"),
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
FOR u IN users
  COLLECT genderGroup = u.gender
   AGGREGATE ages = UNION(u.age, u.age)

  RETURN {
    genderGroup,
    ages,
  }
`, []any{
			map[string]any{"genderGroup": "f", "ages": []any{25, 45, 25, 45}},
			map[string]any{"genderGroup": "m", "ages": []any{31, 36, 69, 31, 36, 69}},
		}, "Should collect and aggregate values by a single key"),
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
  				COLLECT ageGroup = FLOOR(u.age / 5) * 5 
  				AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  				RETURN {
					ageGroup,
    				minAge, 
    				maxAge 
  				}
`, []any{
			map[string]any{"ageGroup": 25, "maxAge": 25, "minAge": 25},
			map[string]any{"ageGroup": 30, "maxAge": 31, "minAge": 31},
			map[string]any{"ageGroup": 35, "maxAge": 36, "minAge": 36},
			map[string]any{"ageGroup": 45, "maxAge": 45, "minAge": 45},
			map[string]any{"ageGroup": 65, "maxAge": 69, "minAge": 69},
		}, "Should aggregate values with calculated grouping"),
		CaseArray(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true,
					salary: 75000,
					department: "IT"
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false,
					salary: 60000,
					department: "Marketing"
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false,
					salary: 80000,
					department: "IT"
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true,
					salary: 95000,
					department: "Management"
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true,
					salary: 70000,
					department: "Marketing"
				}
			]
			FOR u IN users
				COLLECT department = u.department
				AGGREGATE 
					minAge = MIN(u.age), 
					maxAge = MAX(u.age), 
					avgSalary = AVERAGE(u.salary), 
					totalSalary = SUM(u.salary),
					employeeCount = LENGTH(u)
				RETURN {
					department,
					minAge,
					maxAge,
					avgSalary,
					totalSalary,
					employeeCount
				}
		`, []any{
			map[string]any{
				"department":    "IT",
				"minAge":        31,
				"maxAge":        36,
				"avgSalary":     77500.0,
				"totalSalary":   155000,
				"employeeCount": 2,
			},
			map[string]any{
				"department":    "Management",
				"minAge":        69,
				"maxAge":        69,
				"avgSalary":     95000.0,
				"totalSalary":   95000,
				"employeeCount": 1,
			},
			map[string]any{
				"department":    "Marketing",
				"minAge":        25,
				"maxAge":        45,
				"avgSalary":     65000.0,
				"totalSalary":   130000,
				"employeeCount": 2,
			},
		}, "Should aggregate multiple values with complex expressions"),
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
