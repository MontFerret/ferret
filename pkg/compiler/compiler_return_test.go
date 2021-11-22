package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReturn(t *testing.T) {
	Convey("Should compile RETURN NONE", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN NONE
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "null")
	})

	Convey("Should compile RETURN TRUE", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN TRUE
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("Should compile RETURN 1", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 1
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1")
	})

	Convey("Should compile RETURN 1.1", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 1.1
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1.1")
	})

	Convey("Should compile RETURN 'foo'", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 'foo'
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "\"foo\"")
	})

	Convey("Should compile RETURN \"foo\"", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN "foo"
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "\"foo\"")
	})

	Convey("Should compile RETURN \"\"", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN ""
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "\"\"")
	})

	Convey("Should compile RETURN []", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN []
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[]")
	})

	Convey("Should compile RETURN [1, 2, 3, 4]", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN [1, 2, 3, 4]
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[1,2,3,4]")
	})

	Convey("Should compile RETURN ['foo', 'bar', 'qaz']", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN ['foo', 'bar', 'qaz']
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"foo\",\"bar\",\"qaz\"]")
	})

	Convey("Should compile RETURN ['foo', 'bar', 1, 2]", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN ['foo', 'bar', 1, 2]
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "[\"foo\",\"bar\",1,2]")
	})

	Convey("Should compile RETURN {}", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN {}
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{}")
	})

	Convey("Should compile RETURN {a: 'foo', b: 'bar'}", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN { a: "foo", b: "bar" }
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{\"a\":\"foo\",\"b\":\"bar\"}")
	})

	Convey("Should compile RETURN {['a']: 'foo'}", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN { ["a"]: "foo" }
		`)

		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "{\"a\":\"foo\"}")
	})

	SkipConvey("Should compile RETURN (WAITFOR EVENT \"event\" IN obj)", t, func() {
		c := newCompilerWithObservable()

		out, err := c.MustCompile(`
			LET obj = X::VAL("event", ["data"])

			RETURN (WAITFOR EVENT "event" IN obj)
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `"data"`)
	})
}
