package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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
