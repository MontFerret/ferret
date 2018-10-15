package compiler_test

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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
