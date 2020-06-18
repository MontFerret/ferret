package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLet(t *testing.T) {
	Convey("Should compile LET i = NONE RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = NONE
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "null")
	})

	Convey("Should compile LET i = TRUE RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = TRUE
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("Should compile LET i = 1 RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = 1
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1")
	})

	Convey("Should compile LET i = 1.1 RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = 1.1
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1.1")
	})

	Convey("Should compile LET i = 'foo' RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = "foo"
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "\"foo\"")
	})

	Convey("Should compile LET i = [] RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = []
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[]")
	})

	Convey("Should compile LET i = [1, 2, 3] RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = [1, 2, 3]
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3]")
	})

	Convey("Should compile LET i = {} RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = {}
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{}")
	})

	Convey("Should compile LET i = {a: 'foo', b: 1, c: TRUE, d: [], e: {}} RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = {a: 'foo', b: 1, c: TRUE, d: [], e: {}}
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{\"a\":\"foo\",\"b\":1,\"c\":true,\"d\":[],\"e\":{}}")
	})

	Convey("Should compile LET i = (FOR i IN [1,2,3] RETURN i) RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = (FOR i IN [1,2,3] RETURN i)
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3]")
	})

	Convey("Should compile LET i = (FOR i WHILE 0 > 1 RETURN i) RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET i = (FOR i WHILE 0 > 1 RETURN i)
			RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[]")
	})

	Convey("Should compile LET i = { items: [1,2,3]}  FOR el IN i.items RETURN i", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET obj = { items: [1,2,3] }
	
			FOR i IN obj.items
				RETURN i
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3]")
	})

	Convey("Should not compile FOR foo IN foo", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
			FOR foo IN foo
				RETURN foo
		`)

		So(err, ShouldNotBeNil)
	})

	Convey("Should not compile if a variable not defined", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
			RETURN foo
		`)

		So(err, ShouldNotBeNil)
	})

	Convey("Should not compile if a variable is not unique", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
			LET foo = "bar"
			LET foo = "baz"

			RETURN foo
		`)

		So(err, ShouldNotBeNil)
	})
}
