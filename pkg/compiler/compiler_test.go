package compiler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

func TestReturn(t *testing.T) {
	Convey("Should compile RETURN NONE", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN NONE
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "null")
	})

	Convey("Should compile RETURN TRUE", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN TRUE
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("Should compile RETURN 1", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 1
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1")
	})

	Convey("Should compile RETURN 1.1", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 1.1
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1.1")
	})

	Convey("Should compile RETURN 'foo'", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 'foo'
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "\"foo\"")
	})

	Convey("Should compile RETURN \"foo\"", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN "foo"
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "\"foo\"")
	})

	Convey("Should compile RETURN \"\"", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN ""
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "\"\"")
	})

	Convey("Should compile RETURN []", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN []
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[]")
	})

	Convey("Should compile RETURN [1, 2, 3, 4]", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN [1, 2, 3, 4]
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3,4]")
	})

	Convey("Should compile RETURN ['foo', 'bar', 'qaz']", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN ['foo', 'bar', 'qaz']
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"foo\",\"bar\",\"qaz\"]")
	})

	Convey("Should compile RETURN ['foo', 'bar', 1, 2]", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN ['foo', 'bar', 1, 2]
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"foo\",\"bar\",1,2]")
	})

	Convey("Should compile RETURN {}", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN {}
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{}")
	})

	Convey("Should compile RETURN {a: 'foo', b: 'bar'}", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN { a: "foo", b: "bar" }
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{\"a\":\"foo\",\"b\":\"bar\"}")
	})

	Convey("Should compile RETURN {['a']: 'foo'}", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN { ["a"]: "foo" }
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{\"a\":\"foo\"}")
	})
}

func TestFor(t *testing.T) {
	Convey("Should compile FOR i IN [] RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN []
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[]")
	})

	Convey("Should compile FOR i IN [1, 2, 3] RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [1, 2, 3]
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3]")
	})

	Convey("Should compile FOR i, k IN [1, 2, 3] RETURN k", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i, k IN [1, 2, 3]
				RETURN k
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[0,1,2]")
	})

	Convey("Should compile FOR i IN ['foo', 'bar', 'qaz'] RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN ['foo', 'bar', 'qaz']
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"foo\",\"bar\",\"qaz\"]")
	})

	Convey("Should compile FOR i IN {a: 'bar', b: 'foo', c: 'qaz'} RETURN i.name", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN {a: 'bar', b: 'foo', c: 'qaz'}
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		arr := make([]string, 0, 3)
		err = json.Unmarshal(out, &arr)

		So(err, ShouldBeNil)

		sort.Strings(arr)

		out, err = json.Marshal(arr)

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, "[\"bar\",\"foo\",\"qaz\"]")
	})

	Convey("Should compile FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'} RETURN k", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'}
				RETURN k
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		arr := make([]string, 0, 3)
		err = json.Unmarshal(out, &arr)

		So(err, ShouldBeNil)

		sort.Strings(arr)

		out, err = json.Marshal(arr)

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"a\",\"b\",\"c\"]")
	})

	Convey("Should compile FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] RETURN i.name", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}]
				RETURN i.name
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"foo\",\"bar\",\"qaz\"]")
	})

	Convey("Should compile nested FOR operators", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					RETURN {[prop]: val}
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, "[{\"a\":1},{\"a\":2},{\"a\":3}]")
	})

	Convey("Should compile deeply nested FOR operators", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					FOR val2 IN [1, 2, 3]
						RETURN { [prop]: [val, val2] }
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"a":[1,1]},{"a":[1,2]},{"a":[1,3]},{"a":[2,1]},{"a":[2,2]},{"a":[2,3]},{"a":[3,1]},{"a":[3,2]},{"a":[3,3]}]`)
	})

	Convey("Should compile query with a sub query", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR val IN [1, 2, 3]
				RETURN (
					FOR prop IN ["a", "b", "c"]
						RETURN { [prop]: val }
				)
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[[{"a":1},{"b":1},{"c":1}],[{"a":2},{"b":2},{"c":2}],[{"a":3},{"b":3},{"c":3}]]`)
	})

	Convey("Should compile query with variable in a body", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR val IN [1, 2, 3]
				LET sub = (
					FOR prop IN ["a", "b", "c"]
						RETURN { [prop]: val }
				)

				RETURN sub
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[[{"a":1},{"b":1},{"c":1}],[{"a":2},{"b":2},{"c":2}],[{"a":3},{"b":3},{"c":3}]]`)
	})

	Convey("Should compile query with RETURN DISTINCT", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				RETURN DISTINCT i
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[1,2,3,4]`)
	})

	Convey("Should compile query with LIMIT 2", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LIMIT 2
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[1,2]`)
	})

	Convey("Should compile query with LIMIT 2, 2", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LIMIT 2, 2
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[3,4]`)
	})

	Convey("Should compile query with FILTER i > 2", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 2
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[3,4,3]`)
	})

	Convey("Should compile query with FILTER i > 1 AND i < 3", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 1 AND i < 4
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,3,3]`)
	})

	Convey("Should compile query with multiple FILTER statements", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
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
			FOR u IN users
				FILTER u.active == true
				FILTER u.age < 35
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m"},{"active":true,"age":29,"gender":"f"}]`)
	})

	Convey("Should compile query with multiple FILTER statements", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
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
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active == true
				LIMIT 2
				FILTER u.gender == "m"
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m"}]`)
	})

	Convey("Should compile query with SORT statement", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
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
			FOR u IN users
				SORT u.age
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":29,"gender":"f"},{"active":true,"age":31,"gender":"m"},{"active":true,"age":36,"gender":"m"}]`)
	})

	Convey("Should compile query with SORT DESC statement", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
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
			FOR u IN users
				SORT u.age DESC
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":36,"gender":"m"},{"active":true,"age":31,"gender":"m"},{"active":true,"age":29,"gender":"f"}]`)
	})

	Convey("Should compile query with SORT statement with multiple expressions", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
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
			FOR u IN users
				SORT u.age, u.gender
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":29,"gender":"f"},{"active":true,"age":31,"gender":"f"},{"active":true,"age":31,"gender":"m"},{"active":true,"age":36,"gender":"m"}]`)
	})
}

func TestLet(t *testing.T) {
	Convey("Should compile LET i = NONE RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = NONE
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "null")
	})

	Convey("Should compile LET i = TRUE RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = TRUE
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("Should compile LET i = 1 RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = 1
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1")
	})

	Convey("Should compile LET i = 1.1 RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = 1.1
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1.1")
	})

	Convey("Should compile LET i = 'foo' RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = "foo"
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "\"foo\"")
	})

	Convey("Should compile LET i = [] RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = []
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[]")
	})

	Convey("Should compile LET i = [1, 2, 3] RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = [1, 2, 3]
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3]")
	})

	Convey("Should compile LET i = {} RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = {}
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{}")
	})

	Convey("Should compile LET i = {a: 'foo', b: 1, c: TRUE, d: [], e: {}} RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = {a: 'foo', b: 1, c: TRUE, d: [], e: {}}
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{\"a\":\"foo\",\"b\":1,\"c\":true,\"d\":[],\"e\":{}}")
	})

	Convey("Should compile LET i = (FOR i IN [1,2,3] RETURN i) RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET i = (FOR i IN [1,2,3] RETURN i)
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3]")
	})

	Convey("Should compile LET i = { items: [1,2,3]}  FOR el IN i.items RETURN i", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET obj = { items: [1,2,3] }
	
			FOR i IN obj.items
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3]")
	})
}

func TestEqualityOperators(t *testing.T) {
	Convey("Should compile RETURN 2 > 1", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 2 > 1
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})
}

func TestLogicalOperators(t *testing.T) {
	Convey("Should compile RETURN 2 > 1 AND 1 > 0", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 2 > 1 AND 1 > 0
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("Should compile RETURN 2 > 1 OR 1 < 0", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 2 > 1 OR 1 < 0
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("1 || 7  should return 1", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 1 || 7
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1")
	})

	Convey("NONE || 'foo'  should return 'foo'", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN NONE || 'foo'
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `"foo"`)
	})

	Convey("NONE && true  should return null", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN NONE && true
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `null`)
	})

	Convey("'' && true  should return ''", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN '' && true
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `""`)
	})

	Convey("true && 23  should return '23", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN true && 23 
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `23`)
	})

	Convey("NOT TRUE should return false", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN NOT TRUE
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `false`)
	})

	Convey("NOT u.valid should return true", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET u = { valid: false }

			RETURN NOT u.valid
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `true`)
	})
}

func TestMathOperators(t *testing.T) {
	Convey("Integers", t, func() {
		Convey("Should compile RETURN 1 + 1", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 1 + 1
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "2")
		})

		Convey("Should compile RETURN 1 - 1", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 1 - 1
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "0")
		})

		Convey("Should compile RETURN 2*2", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 2*2
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "4")
		})

		Convey("Should compile RETURN 4/2", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 4/2
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "2")
		})

		Convey("Should compile RETURN 5 % 2", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 5 % 2
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "1")
		})
	})

	Convey("Floats", t, func() {
		Convey("Should compile RETURN 1.2 + 1", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 1.2 + 1
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "2.2")
		})

		Convey("Should compile RETURN 1.1 - 1", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 1.1 - 1
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "0.10000000000000009")
		})

		Convey("Should compile RETURN 2.1*2", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 2.1*2
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "4.2")
		})

		Convey("Should compile RETURN 4.4/2", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 4.4/2
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "2.2")
		})

		Convey("Should compile RETURN 5.5 % 2", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 5.5 % 2
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "1")
		})
	})

	Convey("Strings", t, func() {
		Convey("Should concat two strings RETURN 'Foo' + 'Bar'", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN 'Foo' + 'Bar'
		`)

			So(err, ShouldBeNil)
			So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `"FooBar"`)
		})
	})
}

func TestFunctionCall(t *testing.T) {
	Convey("Should compile RETURN TYPENAME(1)", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN TYPENAME(1)
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `"int"`)
	})

	Convey("Should compile WAIT(10) RETURN 1", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			WAIT(10)
			RETURN 1
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `1`)
	})

	Convey("Should compile LET duration = 10 WAIT(duration) RETURN 1", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET duration = 10

			WAIT(duration)

			RETURN 1
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `1`)
	})

	Convey("Should compile function call inside FOR IN statement", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [1, 2, 3, 4]
				LET duration = 10

				WAIT(duration)

				RETURN i * 2
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,4,6,8]`)
	})
}

func TestMember(t *testing.T) {
	Convey("Computed properties", t, func() {
		Convey("Array by literal", func() {
			c := compiler.New()

			prog, err := c.Compile(`
				LET arr = [1,2,3,4]

				RETURN arr[1]
			`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `2`)
		})

		Convey("Array by variable", func() {
			c := compiler.New()

			prog, err := c.Compile(`
				LET arr = [1,2,3,4]
				LET idx = 1

				RETURN arr[idx]
			`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `2`)
		})

		Convey("Object by literal", func() {
			c := compiler.New()

			prog, err := c.Compile(`
				LET obj = { foo: "bar", qaz: "wsx"}

				RETURN obj["qaz"]
			`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"wsx"`)
		})

		Convey("Object by variable", func() {
			c := compiler.New()

			prog, err := c.Compile(`
				LET obj = { foo: "bar", qaz: "wsx"}
				LET key = "qaz"

				RETURN obj[key]
			`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"wsx"`)
		})
	})
}

func TestTernaryOperator(t *testing.T) {
	Convey("Should compile ternary operator", t, func() {
		c := compiler.New()
		prog, err := c.Compile(`
			FOR i IN [1, 2, 3, 4, 5, 6]
				RETURN i < 3 ? i * 3 : i * 2;
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[3,6,6,8,10,12]`)
	})

	Convey("Should compile ternary operator with shortcut", t, func() {
		c := compiler.New()
		prog, err := c.Compile(`
			FOR i IN [1, 2, 3, 4, 5, 6]
				RETURN i < 3 ? : i * 2
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[true,true,6,8,10,12]`)
	})

	Convey("Should compile ternary operator with shortcut with nones", t, func() {
		c := compiler.New()
		prog, err := c.Compile(`
			FOR i IN [NONE, 2, 3, 4, 5, 6]
				RETURN i ? : i
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[null,2,3,4,5,6]`)
	})

	Convey("Should compile ternary operator with default values", t, func() {
		vals := []string{
			"0",
			"0.0",
			"''",
			"NONE",
			"FALSE",
		}

		c := compiler.New()

		for _, val := range vals {
			prog, err := c.Compile(fmt.Sprintf(`
			FOR i IN [%s, 1, 2, 3]
				RETURN i ? i * 2 : 'no value'
		`, val))

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `["no value",2,4,6]`)
		}
	})
}

func TestRangeOperator(t *testing.T) {
	Convey("Should compile RETURN 1..10", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
				RETURN 1..10
			`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[1,2,3,4,5,6,7,8,9,10]`)
	})

	Convey("Should compile FOR i IN 1..10 RETURN i * 2", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
				FOR i IN 1..10
					RETURN i * 2
			`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)
	})

	Convey("Should compile LET arr = 1..10 FOR i IN arr RETURN i * 2", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
				LET arr = 1..10
				FOR i IN arr
					RETURN i * 2
			`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)
	})

	Convey("Should use variables", t, func() {
		out := compiler.New().MustCompile(`
				LET max = 10
				
				FOR i IN 1..max
					RETURN i * 2
		`).MustRun(context.Background())

		So(string(out), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)

		out2 := compiler.New().MustCompile(`
				LET min = 1
				
				FOR i IN min..10
					RETURN i * 2
		`).MustRun(context.Background())

		So(string(out2), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)

		out3 := compiler.New().MustCompile(`
				LET min = 1
				LET max = 10
				
				FOR i IN min..max
					RETURN i * 2
		`).MustRun(context.Background())

		So(string(out3), ShouldEqual, `[2,4,6,8,10,12,14,16,18,20]`)
	})
}

func TestInOperator(t *testing.T) {
	Convey("1 IN [1,2,3] should return true", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 1 IN [1,2,3]
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `true`)
	})

	Convey("4 IN [1,2,3] should return false", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 4 IN [1,2,3]
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `false`)
	})

	Convey("4 NOT IN [1,2,3] should return true", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 4 NOT IN [1,2,3]
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `true`)
	})
}

func TestForTernaryExpression(t *testing.T) {
	Convey("RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6,8,10]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `true`)
	})

	Convey("RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6,8,10]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `[1,2,3,4,5]`)
	})

	Convey("LET res =  foo ? TRUE : (FOR i IN 1..5 RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2) 
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6,8,10]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `true`)
	})

	Convey("LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = FALSE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `[2,4,6,8,10]`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `[1,2,3,4,5]`)
	})
}

func TestParam(t *testing.T) {
	Convey("Should be possible to use as a return value", t, func() {
		out := compiler.New().
			MustCompile(`
			RETURN @param
		`).
			MustRun(context.Background(), runtime.WithParam("param", "foobar"))

		So(string(out), ShouldEqual, `"foobar"`)
	})

	Convey("Should be possible to use as a FOR source", t, func() {
		out := compiler.New().
			MustCompile(`
			FOR i IN @values
			SORT i
			RETURN i
		`).
			MustRun(context.Background(), runtime.WithParam("values", []int{1, 2, 3, 4}))

		So(string(out), ShouldEqual, `[1,2,3,4]`)

		out2 := compiler.New().
			MustCompile(`
			FOR i IN @values
			SORT i
			RETURN i
		`).
			MustRun(context.Background(), runtime.WithParam("values", map[string]int{
				"foo": 1,
				"bar": 2,
				"faz": 3,
				"qaz": 4,
			}))

		So(string(out2), ShouldEqual, `[1,2,3,4]`)
	})

	Convey("Should be possible to use in range", t, func() {
		out := compiler.New().
			MustCompile(`
			FOR i IN @start..@end
			SORT i
			RETURN i
		`).
			MustRun(
				context.Background(),
				runtime.WithParam("start", 1),
				runtime.WithParam("end", 4),
			)

		So(string(out), ShouldEqual, `[1,2,3,4]`)

	})
}

func TestHtml(t *testing.T) {
	Convey("Should load a document", t, func() {
		c := compiler.New()

		out, err := c.MustCompile(`
LET doc = DOCUMENT("http://getbootstrap.com/docs/4.1/components/collapse/", true)

CLICK(doc, "#headingTwo > h5 > button")
WAIT_CLASS(doc, "#collapseTwo", "bar")

RETURN TRUE
			`).Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `"int"`)
	})
}
