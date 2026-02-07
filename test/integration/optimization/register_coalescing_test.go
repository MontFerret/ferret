package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestRegisterCoalescing(t *testing.T) {
	RunUseCases(t, compiler.O1, []UseCase{
		RegistersCase(`
LET a = 10
LET b = a + 1
LET c = b * 2
LET d = c - 3
RETURN d
`, 2, 19),
		RegistersCase(`
LET a = 10
LET b = a + 1
LET c = b * 2
LET d = c - 3
RETURN d
`, 2, 19),
		RegistersCase(`
		LET a = 10
		LET b = a
		LET c = b + 1
		RETURN c
		`, 2, 11),
		RegistersCase(`
		LET a = 1
		LET b = a
		LET c = b
		RETURN c
		`, 1, 1),
		RegistersCase(`
		LET a = 1
		LET b = 2
		LET c = a + b
		RETURN c
		`, 2, 3),
		RegistersArrayCase(`
		LET a = 10
		LET arr = [a, a + 1, a + 2, a + 3]
		RETURN arr
		`, 3, []any{10, 11, 12, 13}, "Flat array literal with expression elems"),
		RegistersObjectCase(`
		LET x = 5
		LET obj = {
		 a: x,
		 b: x + 1,
		 c: (x + 1) * 2,
		 d: (x + 2) * 3
		}
		RETURN obj
		`, 4, map[string]any{
			"a": 5,
			"b": 6,
			"c": 12,
			"d": 21,
		}, "Flat object literal with expression elems"),
		RegistersArrayCase(`
		LET base = 100
		LET items = [ {a: base}, {a: 2 }, {a: 3} ]
		RETURN items
		`, 4, []any{
			map[string]any{"a": 100},
			map[string]any{"a": 2},
			map[string]any{"a": 3},
		}, "Nested literals (array of objects)"),
		RegistersObjectCase(`
		LET x = 10
		LET doc = { meta:{ a: x }, data:[x, 3], sum:(x*2)+(x*3) }
		RETURN doc
		`, 4, map[string]any{
			"meta": map[string]any{"a": 10},
			"data": []any{10, 3},
			"sum":  50,
		}, "Object containing arrays + nested computed values"),
		RegistersArrayCase(`
LET a = [10,20,30,40]
LET i = 1
LET out = [a[i], a[i+1], a[i+2]]
RETURN out`, 4, []any{20, 30, 40}, "Computed index pattern"),
		RegistersObjectCase(`
LET k="price" 
LET v=123
LET obj = { [k]: v, ["qty"]:2, ["total"]: v*2 }
RETURN obj`, 4, map[string]any{
			"price": 123,
			"qty":   2,
			"total": 246,
		}, "Computed keys in object literal"),
		RegistersArrayCase(`
LET o = { a: 1 }
LET arr = [o, o, o]
RETURN arr`, 3, []any{
			map[string]any{"a": 1},
			map[string]any{"a": 1},
			map[string]any{"a": 1},
		}, "Alias pitfall (same object ref in array)"),
		RegistersCase(`
RETURN FIRST([])?.foo
`, 2, nil, "Optional chaining with FIRST on empty array"),
		RegistersArrayCase(`
FOR x IN [1,2,3,4,5]
  LET row = { x:x, y:x*2, z:x*3 }
  RETURN row
`, 6, []any{
			map[string]any{"x": 1, "y": 2, "z": 3},
			map[string]any{"x": 2, "y": 4, "z": 6},
			map[string]any{"x": 3, "y": 6, "z": 9},
			map[string]any{"x": 4, "y": 8, "z": 12},
			map[string]any{"x": 5, "y": 10, "z": 15},
		}, "Register reuse across loop iterations"),
		RegistersCase(`
LET x = 1
RETURN (x > 0) ? (x + 1) : (x - 1)
`, 2, 2, "Simple ternary with arithmetic on both sides"),
		RegistersCase(`
LET a = 1
LET b = 2
RETURN (a == b)
  ? ((a + b) * 2)
  : ((a - b) / 2)
`, 3, -0.5, "Equality + ternary nesting (diamond inside diamond)"),
		RegistersCase(`
LET x = 1
LET y = 2
RETURN (x > 0 AND y > 0) ? (x + y) : (x - y)
`, 3, 3, "Boolean combos with short-circuit"),
		RegistersCase(`
LET x = 2
RETURN ((x + 1) * (x + 2) == (x + 3) * (x + 4))
  ? (x * x + 1)
  : (x * x - 1)
`, 5, 3, "Chained comparisons + math cascade"),
		RegistersCase(`
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
		`, 13, []any{
			map[string]any{
				"department": "IT",
				"budget":     500000.0,
				"stats": map[string]any{
					"avgAge":      33.5,
					"totalSalary": 155000.0,
					"count":       2.0,
				},
			},
			map[string]any{
				"department": "Marketing",
				"budget":     300000.0,
				"stats": map[string]any{
					"avgAge":      35.0,
					"totalSalary": 130000.0,
					"count":       2.0,
				},
			},
			map[string]any{
				"department": "Management",
				"budget":     200000.0,
				"stats": map[string]any{
					"avgAge":      69.0,
					"totalSalary": 95000.0,
					"count":       1.0,
				},
			},
		}, "Should handle nested FOR loops with COLLECT AGGREGATE"),
		RegistersCase(`RETURN FIRST([])?.foo`, 2, nil, "Optional chaining with array access and property access"),
		RegistersArrayCase(`
			LET users = [
				{ gender: "m", age: null },
				{ gender: "m", age: 40 },
				{ gender: "f", age: 20 },
				{ gender: "f", age: null }
			]
			FOR u IN users
				COLLECT gender = u.gender
				AGGREGATE 
					count = COUNT(u.age), 
					sum = SUM(u.age), 
					avg = AVERAGE(u.age)
				RETURN {
					gender, count, sum, avg
				}
		`, 9, []any{
			map[string]any{"gender": "f", "count": 2, "sum": 20, "avg": 20},
			map[string]any{"gender": "m", "count": 2, "sum": 40, "avg": 40},
		}, "Should skip nulls in COUNT, SUM, AVG"),
		RegistersArrayCase(
			`FOR i IN [ 1, 2, 3, 4, 1, 3 ]
							RETURN DISTINCT i
		`,
			3, []any{1, 2, 3, 4},
		),
	})
}
