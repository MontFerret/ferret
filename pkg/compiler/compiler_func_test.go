package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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
