package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForStepNested(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			LET props = ["a"]
			LET counter1 = 1
			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET prop = props[i]
				LET vals = [1, 2, 3]
				LET counter2 = 3
				FOR j = 0 WHILE j < counter2 STEP j = j + 1
					LET val = vals[j]
					RETURN {[prop]: val}
`, []any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(`
			FOR val IN 1..3
				LET props = ["a"]
				LET counter1 = 1
				FOR j = 0 WHILE j < counter1 STEP j = j + 1
					LET prop = props[j]
					RETURN {[prop]: val}
`, []any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(`
			LET props = ["a"]
			LET counter1 = 1
			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET prop = props[i]
				FOR val IN 1..3
					RETURN {[prop]: val}
`, []any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(`
			LET props = ["a"]
			LET counter1 = 1
			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET prop = props[i]
				LET vals = [1, 2, 3]
				LET counter2 = 3
				FOR j = 0 WHILE j < counter2 STEP j = j + 1
					LET val = vals[j]
					LET vals2 = [1, 2, 3]
					LET counter3 = 3
					FOR k = 0 WHILE k < counter3 STEP k = k + 1
						LET val2 = vals2[k]
						RETURN { [prop]: [val, val2] }
`, []any{map[string]any{"a": []int{1, 1}}, map[string]any{"a": []int{1, 2}}, map[string]any{"a": []int{1, 3}}, map[string]any{"a": []int{2, 1}}, map[string]any{"a": []int{2, 2}}, map[string]any{"a": []int{2, 3}}, map[string]any{"a": []int{3, 1}}, map[string]any{"a": []int{3, 2}}, map[string]any{"a": []int{3, 3}}},
		),
		CaseArray(`
			LET vals = [1, 2, 3]
			LET counter1 = 3
			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET val = vals[i]
				LET props = ["a", "b", "c"]
				LET counter2 = 3
				RETURN (
					FOR j = 0 WHILE j < counter2 STEP j = j + 1
						LET prop = props[j]
						RETURN { [prop]: val }
				)
`, []any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
		),
		CaseArray(`
			LET vals = [1, 2, 3]
			LET counter1 = 3
			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET val = vals[i]
				LET props = ["a", "b", "c"]
				LET counter2 = 3
				LET sub = (
					FOR j = 0 WHILE j < counter2 STEP j = j + 1
						LET prop = props[j]
						RETURN { [prop]: val }
				)

				RETURN sub
`, []any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}}),
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
				}
			]

			LET targetSkills = ["JavaScript", "Python", "Go", "Java"]
			LET counter1 = 3
			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET u = users[i]
				LET counter2 = 4
				LET matchingSkills = (
					FOR j = 0 WHILE j < counter2 STEP j = j + 1
						LET skill = targetSkills[j]
						FILTER skill IN u.skills
						RETURN skill
				)

				RETURN {
					name: u.name,
					matchingSkills: matchingSkills,
					matchCount: LENGTH(matchingSkills),
					hasJavaScript: "JavaScript" IN u.skills,
					hasPython: "Python" IN u.skills
				}
		`, []any{
			map[string]any{
				"name":           "John",
				"matchingSkills": []any{"JavaScript", "Python", "Go"},
				"matchCount":     3,
				"hasJavaScript":  true,
				"hasPython":      true,
			},
			map[string]any{
				"name":           "Jane",
				"matchingSkills": []any{"Python", "Java"},
				"matchCount":     2,
				"hasJavaScript":  false,
				"hasPython":      true,
			},
			map[string]any{
				"name":           "Bob",
				"matchingSkills": []any{"Go"},
				"matchCount":     1,
				"hasJavaScript":  false,
				"hasPython":      false,
			},
		}, "Should handle nested FOR loops with array operations"),
		CaseArray(`
			LET departments = ["IT", "Marketing", "HR"]
			LET budgets = [1000000, 500000, 300000]
			LET headcounts = [20, 10, 5]
			LET counter1 = 3

			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET dept = departments[i]
				LET budget = budgets[i]
				LET headcount = headcounts[i]

				FOR j IN 1..3
					LET allocation = (
						j == 1 ? 0.5 :
						j == 2 ? 0.3 :
						0.2
					)

					LET category = (
						j == 1 ? "Salaries" :
						j == 2 ? "Equipment" :
						"Miscellaneous"
					)

					RETURN {
						department: dept,
						category: category,
						allocation: allocation,
						amount: budget * allocation,
						headcount: headcount
					}
		`, []any{
			map[string]any{
				"allocation": 0.3,
				"amount":     300000,
				"category":   "Equipment",
				"department": "IT",
				"headcount":  20,
			},
			map[string]any{
				"allocation": 0.3,
				"amount":     300000,
				"category":   "Equipment",
				"department": "IT",
				"headcount":  20,
			},
			map[string]any{
				"allocation": 0.2,
				"amount":     200000,
				"category":   "Miscellaneous",
				"department": "IT",
				"headcount":  20,
			},
			map[string]any{
				"allocation": 0.3,
				"amount":     150000,
				"category":   "Equipment",
				"department": "Marketing",
				"headcount":  10,
			},
			map[string]any{
				"allocation": 0.3,
				"amount":     150000,
				"category":   "Equipment",
				"department": "Marketing",
				"headcount":  10,
			},
			map[string]any{
				"allocation": 0.2,
				"amount":     100000,
				"category":   "Miscellaneous",
				"department": "Marketing",
				"headcount":  10,
			},
			map[string]any{
				"allocation": 0.3,
				"amount":     90000,
				"category":   "Equipment",
				"department": "HR",
				"headcount":  5,
			},
			map[string]any{
				"allocation": 0.3,
				"amount":     90000,
				"category":   "Equipment",
				"department": "HR",
				"headcount":  5,
			},
			map[string]any{
				"allocation": 0.2,
				"amount":     60000,
				"category":   "Miscellaneous",
				"department": "HR",
				"headcount":  5,
			},
		}, "Should handle nested FOR loops with conditional logic"),
		CaseArray(`
			LET users = [
				{
					id: 1,
					name: "John",
					department: "IT",
					projects: [
						{ id: "p1", name: "Project A", status: "active" },
						{ id: "p2", name: "Project B", status: "completed" }
					]
				},
				{
					id: 2,
					name: "Jane",
					department: "Marketing",
					projects: [
						{ id: "p3", name: "Project C", status: "active" },
						{ id: "p4", name: "Project D", status: "active" }
					]
				},
				{
					id: 3,
					name: "Bob",
					department: "IT",
					projects: [
						{ id: "p5", name: "Project E", status: "pending" }
					]
				}
			]
			
			LET end = LENGTH(users) - 1
			LET counter1 = (FOR i IN 0..end RETURN i)

			FOR i IN 0..end
				LET u = users[i]
				LET end2 = LENGTH(u.projects) - 1
				LET counter2 = (FOR i IN 0..end2 RETURN i)
				LET activeProjects = (
					FOR j IN 0..end2
						LET p = u.projects[j]
						FILTER p.status == "active"
						RETURN p
				)

				FILTER LENGTH(activeProjects) > 0

				LET end3 = LENGTH(activeProjects) - 1
				LET counter3 = (FOR i IN 0..end3 RETURN i)

				RETURN {
					user: u.name,
					department: u.department,
					activeProjects: (
						FOR k WHILE W_POP(counter3)
							LET p = activeProjects[k]
							RETURN p.name
					)
				}
		`, []any{
			map[string]any{
				"user":           "John",
				"department":     "IT",
				"activeProjects": []any{"Project A"},
			},
			map[string]any{
				"user":           "Jane",
				"department":     "Marketing",
				"activeProjects": []any{"Project C", "Project D"},
			},
		}, "Should handle nested FOR loops with complex data transformation"),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]

			FOR n IN 0..1
				LET counter1 = 4
				FOR i = 0 WHILE i < counter1 STEP i = i + 1
					LET s = strs[i]
					SORT s
					RETURN CONCAT(s, n)
`, []any{"abc0", "bar0", "foo0", "qaz0", "abc1", "bar1", "foo1", "qaz1"}),
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

			FOR n IN 0..1
				LET counter1 = 3
				FOR i = 0 WHILE i < counter1 STEP i = i + 1
					LET u = users[i]
					SORT u.gender, u.age
					RETURN CONCAT(u.name, n)
`, []any{"Angela0", "Ron0", "Bob0", "Angela1", "Ron1", "Bob1"}),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]
			LET counter1 = [0,0]

			FOR n WHILE W_POP(counter1)
				FOR m IN 0..1
					LET counter2 = 4
					FOR i = 0 WHILE i < counter2 STEP i = i + 1
						LET s = strs[i]
						SORT s
						RETURN CONCAT(s, n, m)
`, []any{"abc00", "bar00", "foo00", "qaz00", "abc01", "bar01", "foo01", "qaz01", "abc10", "bar10", "foo10", "qaz10", "abc11", "bar11", "foo11", "qaz11"}),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]
			LET counter1 = 2

			FOR n = 0 WHILE n < counter1 STEP n = n + 1
				LET counter2 = [0,0,0,0]
				FOR i WHILE W_POP(counter2)
					LET s = strs[i]
					SORT s
					FOR m IN 0..1
						RETURN CONCAT(s, n, m)
`, []any{"abc00", "abc01", "bar00", "bar01", "foo00", "foo01", "qaz00", "qaz01", "abc10", "abc11", "bar10", "bar11", "foo10", "foo11", "qaz10", "qaz11"}),
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
			FOR n IN 0..1
				LET counter1 = 5
				FOR i = 0 WHILE i < counter1 STEP i = i + 1
					LET u = users[i]
					COLLECT gender = u.gender
					RETURN CONCAT(gender, n)
`, []any{"f0", "m0", "f1", "m1"}),
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
			LET counter1 = 5
			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET u = users[i]
				COLLECT gender = u.gender
				LET counter2 = [0,0]
				FOR n WHILE W_POP(counter2)
					RETURN CONCAT(gender, n)
`, []any{"f0", "f1", "m0", "m1"}),
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
			FOR n IN 0..1
				LET counter1 = 5
				FOR i = 0 WHILE i < counter1 STEP i = i + 1
					LET u = users[i]
					COLLECT gender = u.gender INTO genders
					RETURN {
						gender: CONCAT(gender, n),
						values: genders
					}
`, []any{
			map[string]any{
				"gender": "f0",
				"values": []any{
					map[string]any{
						"i": 1,
						"u": map[string]any{
							"active":  true,
							"age":     25,
							"gender":  "f",
							"married": false,
						},
					},
					map[string]any{
						"i": 4,
						"u": map[string]any{
							"active":  true,
							"age":     45,
							"gender":  "f",
							"married": true,
						},
					},
				},
			},
			map[string]any{
				"gender": "m0",
				"values": []any{
					map[string]any{
						"i": 0,
						"u": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					map[string]any{
						"i": 2,
						"u": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					map[string]any{
						"i": 3,
						"u": map[string]any{
							"active":  false,
							"age":     69,
							"gender":  "m",
							"married": true,
						},
					},
				},
			},
			map[string]any{
				"gender": "f1",
				"values": []any{
					map[string]any{
						"i": 1,
						"u": map[string]any{
							"active":  true,
							"age":     25,
							"gender":  "f",
							"married": false,
						},
					},
					map[string]any{
						"i": 4,
						"u": map[string]any{
							"active":  true,
							"age":     45,
							"gender":  "f",
							"married": true,
						},
					},
				},
			},
			map[string]any{
				"gender": "m1",
				"values": []any{
					map[string]any{
						"i": 0,
						"u": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					map[string]any{
						"i": 2,
						"u": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					map[string]any{
						"i": 3,
						"u": map[string]any{
							"active":  false,
							"age":     69,
							"gender":  "m",
							"married": true,
						},
					},
				},
			},
		}),
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
			LET counter1 = 5
			FOR i = 0 WHILE i < counter1 STEP i = i + 1
				LET u = users[i]
				COLLECT gender = u.gender INTO genders
				FOR n IN 0..1
					RETURN {
						gender: CONCAT(gender, n),
						values: genders
					}
`, []any{
			map[string]any{
				"gender": "f0",
				"values": []any{
					map[string]any{
						"i": 1,
						"u": map[string]any{
							"active":  true,
							"age":     25,
							"gender":  "f",
							"married": false,
						},
					},
					map[string]any{
						"i": 4,
						"u": map[string]any{
							"active":  true,
							"age":     45,
							"gender":  "f",
							"married": true,
						},
					},
				},
			},
			map[string]any{
				"gender": "f1",
				"values": []any{
					map[string]any{
						"i": 1,
						"u": map[string]any{
							"active":  true,
							"age":     25,
							"gender":  "f",
							"married": false,
						},
					},
					map[string]any{
						"i": 4,
						"u": map[string]any{
							"active":  true,
							"age":     45,
							"gender":  "f",
							"married": true,
						},
					},
				},
			},
			map[string]any{
				"gender": "m0",
				"values": []any{
					map[string]any{
						"i": 0,
						"u": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					map[string]any{
						"i": 2,
						"u": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					map[string]any{
						"i": 3,
						"u": map[string]any{
							"active":  false,
							"age":     69,
							"gender":  "m",
							"married": true,
						},
					},
				},
			},
			map[string]any{
				"gender": "m1",
				"values": []any{
					map[string]any{
						"i": 0,
						"u": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					map[string]any{
						"i": 2,
						"u": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					map[string]any{
						"i": 3,
						"u": map[string]any{
							"active":  false,
							"age":     69,
							"gender":  "m",
							"married": true,
						},
					},
				},
			},
		}),
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
			FOR n IN 0..1
				LET counter1 = 5
				FOR i = 0 WHILE i < counter1 STEP i = i + 1
				  LET u = users[i]
				  COLLECT gender = u.gender
					AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)

				  RETURN {
					gender: CONCAT(gender, n),
					minAge,
					maxAge
				  }
`, []any{
			map[string]any{
				"gender": "f0",
				"maxAge": 45,
				"minAge": 25,
			},
			map[string]any{
				"gender": "m0",
				"maxAge": 69,
				"minAge": 31,
			},
			map[string]any{
				"gender": "f1",
				"maxAge": 45,
				"minAge": 25,
			},
			map[string]any{
				"gender": "m1",
				"maxAge": 69,
				"minAge": 31,
			},
		}),
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

			FOR i = 0 WHILE i < LENGTH(users) STEP i = i + 1
			  LET u = users[i]
			  COLLECT gender = u.gender
				AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
				FOR n IN 0..1
				  RETURN {
					gender: CONCAT(gender, n),
					minAge,
					maxAge
				  }
`, []any{
			map[string]any{
				"gender": "f0",
				"maxAge": 45,
				"minAge": 25,
			},
			map[string]any{
				"gender": "f1",
				"maxAge": 45,
				"minAge": 25,
			},
			map[string]any{
				"gender": "m0",
				"maxAge": 69,
				"minAge": 31,
			},
			map[string]any{
				"gender": "m1",
				"maxAge": 69,
				"minAge": 31,
			},
		}),
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
			FOR n IN 0..1
				FOR i = 0 WHILE i < LENGTH(users) STEP i = i + 1
  				  LET u = users[i]
  				  COLLECT AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)

				  RETURN {
					iteration: n,
					minAge,
					maxAge
				  }
`, []any{
			map[string]any{
				"iteration": 0,
				"maxAge":    69,
				"minAge":    25,
			},
			map[string]any{
				"iteration": 1,
				"maxAge":    69,
				"minAge":    25,
			},
		}),
		CaseArray(`
			LET departments = [
				{ name: "IT", budget: 500000 },
				{ name: "Marketing", budget: 300000 },
				{ name: "Management", budget: 200000 }
			]
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
			FOR i WHILE UNTIL(LENGTH(departments))
				LET d = departments[i]
				LET deptUsers = (
					FOR j = 0 WHILE j < LENGTH(users) STEP j = j + 1
						LET u = users[j]
						FILTER u.department == d.name
						RETURN u
				)
				LET stats = (
					FOR k = 0 WHILE k < LENGTH(deptUsers) STEP k = k + 1
						LET u = deptUsers[k]
						COLLECT AGGREGATE 
							avgAge = AVERAGE(u.age),
							totalSalary = SUM(u.salary),
							kount = LENGTH(u)
						RETURN {
							avgAge,
							totalSalary,
							count: kount
						}
				)
				RETURN {
					department: d.name,
					budget: d.budget,
					stats: stats[0]
				}
		`, []any{
			map[string]any{
				"department": "IT",
				"budget":     500000,
				"stats": map[string]any{
					"avgAge":      33.5,
					"totalSalary": 155000,
					"count":       2,
				},
			},
			map[string]any{
				"department": "Marketing",
				"budget":     300000,
				"stats": map[string]any{
					"avgAge":      35.0,
					"totalSalary": 130000,
					"count":       2,
				},
			},
			map[string]any{
				"department": "Management",
				"budget":     200000,
				"stats": map[string]any{
					"avgAge":      69.0,
					"totalSalary": 95000,
					"count":       1,
				},
			},
		}, "Should handle nested FOR loops with COLLECT AGGREGATE"),
	}, vm.WithFunctions(ForWhileHelpers()))
}
