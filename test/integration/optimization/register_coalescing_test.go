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
	})
}
