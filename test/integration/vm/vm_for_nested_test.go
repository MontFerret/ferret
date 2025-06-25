package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestForNested(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					RETURN {[prop]: val}
`, []any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(`
			FOR val IN 1..3
				FOR prop IN ["a"]
					RETURN {[prop]: val}
`, []any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(`
			FOR prop IN ["a"]
				FOR val IN 1..3
					RETURN {[prop]: val}
`, []any{map[string]any{"a": 1}, map[string]any{"a": 2}, map[string]any{"a": 3}},
		),
		CaseArray(`
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					FOR val2 IN [1, 2, 3]
						RETURN { [prop]: [val, val2] }
`, []any{map[string]any{"a": []int{1, 1}}, map[string]any{"a": []int{1, 2}}, map[string]any{"a": []int{1, 3}}, map[string]any{"a": []int{2, 1}}, map[string]any{"a": []int{2, 2}}, map[string]any{"a": []int{2, 3}}, map[string]any{"a": []int{3, 1}}, map[string]any{"a": []int{3, 2}}, map[string]any{"a": []int{3, 3}}},
		),
		CaseArray(`
			FOR val IN [1, 2, 3]
				RETURN (
					FOR prop IN ["a", "b", "c"]
						RETURN { [prop]: val }
				)
`, []any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
		),
		CaseArray(`
			FOR val IN [1, 2, 3]
			LET sub = (
				FOR prop IN ["a", "b", "c"]
					RETURN { [prop]: val }
			)

			RETURN sub
`, []any{[]any{map[string]any{"a": 1}, map[string]any{"b": 1}, map[string]any{"c": 1}}, []any{map[string]any{"a": 2}, map[string]any{"b": 2}, map[string]any{"c": 2}}, []any{map[string]any{"a": 3}, map[string]any{"b": 3}, map[string]any{"c": 3}}},
		),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]
			
			FOR s IN strs
				SORT s
				FOR n IN 0..1
					RETURN CONCAT(s, n)
`, []any{"abc0", "abc1", "bar0", "bar1", "foo0", "foo1", "qaz0", "qaz1"}),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]
			
			FOR n IN 0..1
				FOR s IN strs
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
				FOR u IN users
					SORT u.gender, u.age
					RETURN CONCAT(u.name, n)
`, []any{"Angela0", "Ron0", "Bob0", "Angela1", "Ron1", "Bob1"}),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]
			
			FOR n IN 0..1
				FOR m IN 0..1
					FOR s IN strs
						SORT s
						RETURN CONCAT(s, n, m)
`, []any{"abc00", "bar00", "foo00", "qaz00", "abc01", "bar01", "foo01", "qaz01", "abc10", "bar10", "foo10", "qaz10", "abc11", "bar11", "foo11", "qaz11"}),
		CaseArray(`
			LET strs = ["foo", "bar", "qaz", "abc"]
			
			FOR n IN 0..1
				FOR s IN strs
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
				FOR i IN users
					COLLECT gender = i.gender
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
			FOR i IN users
				COLLECT gender = i.gender
				FOR n IN 0..1
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
				FOR i IN users
					COLLECT gender = i.gender INTO genders
					RETURN {
						gender: CONCAT(gender, n),
						values: genders
					}
`, []any{map[string]any{
			"gender": "f0",
			"values": []map[string]any{
				{
					"i": map[string]any{
						"active":  true,
						"age":     25,
						"gender":  "f",
						"married": false,
					},
				},
				{
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
				"gender": "m0",
				"values": []map[string]any{
					{
						"i": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					{
						"i": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					{
						"i": map[string]any{
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
				"values": []map[string]any{
					{
						"i": map[string]any{
							"active":  true,
							"age":     25,
							"gender":  "f",
							"married": false,
						},
					},
					{
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
				"gender": "m1",
				"values": []map[string]any{
					{
						"i": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					{
						"i": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					{
						"i": map[string]any{
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
			FOR i IN users
				COLLECT gender = i.gender INTO genders
				FOR n IN 0..1
					RETURN {
						gender: CONCAT(gender, n),
						values: genders
					}
`, []any{map[string]any{
			"gender": "f0",
			"values": []map[string]any{
				{
					"i": map[string]any{
						"active":  true,
						"age":     25,
						"gender":  "f",
						"married": false,
					},
				},
				{
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
				"gender": "f1",
				"values": []map[string]any{
					{
						"i": map[string]any{
							"active":  true,
							"age":     25,
							"gender":  "f",
							"married": false,
						},
					},
					{
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
				"gender": "m0",
				"values": []map[string]any{
					{
						"i": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					{
						"i": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					{
						"i": map[string]any{
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
				"values": []map[string]any{
					{
						"i": map[string]any{
							"active":  true,
							"age":     31,
							"gender":  "m",
							"married": true,
						},
					},
					{
						"i": map[string]any{
							"active":  true,
							"age":     36,
							"gender":  "m",
							"married": false,
						},
					},
					{
						"i": map[string]any{
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
				FOR u IN users
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
			
			FOR u IN users
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
				FOR u IN users
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
	})
}
