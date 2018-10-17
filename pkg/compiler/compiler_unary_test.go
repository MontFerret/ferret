package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUnaryOperator(t *testing.T) {
	Convey("RETURN !{BOOLEAN}", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			RETURN !TRUE
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `false`)

		out2, err := c.MustCompile(`
			RETURN !FALSE
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `true`)
	})

	Convey("RETURN foo ? TRUE : FALSE ", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET foo = TRUE
			RETURN foo ? TRUE : FALSE
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `true`)

		out2, err := c.MustCompile(`
			LET foo = TRUE
			RETURN !foo ? TRUE : FALSE
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `false`)
	})

	Convey("RETURN { enabled: !val}", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET val = ""
			RETURN { enabled: !val }
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `{"enabled":true}`)

		out2, err := c.MustCompile(`
			LET val = ""
			RETURN { enabled: !!val }
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `{"enabled":false}`)
	})

	Convey("RETURN -v", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET v = 1
			RETURN -v
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `-1`)
	})

	Convey("RETURN +v", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET v = -1
			RETURN +v
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `-1`)
	})
}
