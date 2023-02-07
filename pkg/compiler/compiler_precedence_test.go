package compiler_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPrecedence(t *testing.T) {
	Convey("Math operators", t, func() {
		Convey("2 + 2 * 2", func() {
			c := compiler.New()

			p := c.MustCompile(`RETURN 2 + 2 * 2`)

			out := p.MustRun(context.Background())

			So(string(out), ShouldEqual, "6")
		})

		Convey("2 * 2 + 2", func() {
			c := compiler.New()

			p := c.MustCompile(`RETURN 2 * 2 + 2`)

			out := p.MustRun(context.Background())

			So(string(out), ShouldEqual, "6")
		})

		Convey("2 * (2 + 2)", func() {
			c := compiler.New()

			p := c.MustCompile(`RETURN 2 * (2 + 2)`)

			out := p.MustRun(context.Background())

			So(string(out), ShouldEqual, "8")
		})
	})

	Convey("Logical", t, func() {
		Convey("TRUE OR TRUE AND FALSE", func() {
			c := compiler.New()

			p := c.MustCompile(`RETURN TRUE OR TRUE AND FALSE`)

			out := p.MustRun(context.Background())

			So(string(out), ShouldEqual, "true")
		})

		Convey("FALSE AND TRUE OR TRUE", func() {
			c := compiler.New()

			p := c.MustCompile(`RETURN FALSE AND TRUE OR TRUE`)

			out := p.MustRun(context.Background())

			So(string(out), ShouldEqual, "true")
		})
	})
}
