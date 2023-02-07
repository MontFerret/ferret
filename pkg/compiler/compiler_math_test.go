package compiler_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
)

func TestMathOperators(t *testing.T) {
	Convey("Integers", t, func() {
		Convey("Should compile RETURN 1 + 1", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 1 + 1
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "2")
		})

		Convey("Should compile RETURN 1 - 1", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 1 - 1
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "0")
		})

		Convey("Should compile RETURN 2*2", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 2*2
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "4")
		})

		Convey("Should compile RETURN 4/2", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 4/2
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "2")
		})

		Convey("Should compile RETURN 5 % 2", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 5 % 2
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "1")
		})
	})

	Convey("Floats", t, func() {
		Convey("Should compile RETURN 1.2 + 1", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 1.2 + 1
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "2.2")
		})

		Convey("Should compile RETURN 1.1 - 1", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 1.1 - 1
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "0.10000000000000009")
		})

		Convey("Should compile RETURN 2.1*2", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 2.1*2
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "4.2")
		})

		Convey("Should compile RETURN 4.4/2", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 4.4/2
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "2.2")
		})

		Convey("Should compile RETURN 5.5 % 2", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 5.5 % 2
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, "1")
		})
	})

	Convey("Strings", t, func() {
		Convey("Should concat two strings RETURN 'Foo' + 'Bar'", func() {
			c := compiler.New()

			p, err := c.Compile(`
			RETURN 'Foo' + 'Bar'
		`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `"FooBar"`)
		})
	})
}
