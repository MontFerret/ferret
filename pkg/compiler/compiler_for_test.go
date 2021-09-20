package compiler_test

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

func TestFor(t *testing.T) {
	Convey("Should compile FOR i IN [] RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN []
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[]")
	})

	Convey("Should compile FOR i IN [1, 2, 3] RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN [1, 2, 3]
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3]")
	})

	Convey("Should not allocate memory if NONE is a return statement", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN 0..100
				RETURN NONE
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[]")
	})

	Convey("Should compile FOR i, k IN [1, 2, 3] RETURN k", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i, k IN [1, 2, 3]
				RETURN k
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[0,1,2]")
	})

	Convey("Should compile FOR i IN ['foo', 'bar', 'qaz'] RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN ['foo', 'bar', 'qaz']
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"foo\",\"bar\",\"qaz\"]")
	})

	Convey("Should compile FOR i IN {a: 'bar', b: 'foo', c: 'qaz'} RETURN i.name", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN {a: 'bar', b: 'foo', c: 'qaz'}
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

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

		p, err := c.Compile(`
			FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'}
				RETURN k
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

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

		p, err := c.Compile(`
			FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}]
				RETURN i.name
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"foo\",\"bar\",\"qaz\"]")
	})

	Convey("Should compile nested FOR operators", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					RETURN {[prop]: val}
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, "[{\"a\":1},{\"a\":2},{\"a\":3}]")
	})

	Convey("Should compile deeply nested FOR operators", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					FOR val2 IN [1, 2, 3]
						RETURN { [prop]: [val, val2] }
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"a":[1,1]},{"a":[1,2]},{"a":[1,3]},{"a":[2,1]},{"a":[2,2]},{"a":[2,3]},{"a":[3,1]},{"a":[3,2]},{"a":[3,3]}]`)
	})

	Convey("Should compile query with a sub query", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR val IN [1, 2, 3]
				RETURN (
					FOR prop IN ["a", "b", "c"]
						RETURN { [prop]: val }
				)
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[[{"a":1},{"b":1},{"c":1}],[{"a":2},{"b":2},{"c":2}],[{"a":3},{"b":3},{"c":3}]]`)
	})

	Convey("Should compile query with variable in a body", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR val IN [1, 2, 3]
				LET sub = (
					FOR prop IN ["a", "b", "c"]
						RETURN { [prop]: val }
				)

				RETURN sub
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[[{"a":1},{"b":1},{"c":1}],[{"a":2},{"b":2},{"c":2}],[{"a":3},{"b":3},{"c":3}]]`)
	})

	Convey("Should compile query with RETURN DISTINCT", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				RETURN DISTINCT i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[1,2,3,4]`)
	})
}

func BenchmarkForEmpty(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR i IN []
				RETURN i
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForArray(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR i IN [1,2,3]
				RETURN i
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForObject(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR i IN {a: 'bar', b: 'foo', c: 'qaz'}
				RETURN i
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForNested(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					RETURN {[prop]: val}
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForNested2(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR prop IN ["a"]
				FOR val IN [1, 2, 3]
					FOR val2 IN ["b"]
						RETURN { [prop]: [val, val2] }
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForSub(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR val IN [1, 2, 3]
				RETURN (
					FOR prop IN ["a", "b", "c"]
						RETURN { [prop]: val }
				)
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForSub2(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR val IN [1, 2, 3]
				LET sub = (
					FOR prop IN ["a", "b", "c"]
						RETURN { [prop]: val }
				)

				RETURN sub
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForDistinct(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				RETURN DISTINCT i
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForLimit(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT 2
				RETURN i
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForLimitOffset(b *testing.B) {
	p := compiler.New().MustCompile(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT 4, 2
				RETURN i
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForSort(b *testing.B) {
	p := compiler.New().MustCompile(`
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

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForSort2(b *testing.B) {
	p := compiler.New().MustCompile(`
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
				SORT u.age, u.gender
				RETURN u
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkForSortDesc(b *testing.B) {
	p := compiler.New().MustCompile(`
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

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
