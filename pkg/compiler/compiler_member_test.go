package compiler_test

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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

		Convey("Object by literal with property defined as a string", func() {
			c := compiler.New()

			prog, err := c.Compile(`
				LET obj = { "foo": "bar", "qaz": "wsx"}

				RETURN obj["qaz"]
			`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"wsx"`)
		})

		Convey("Object by literal with property defined as a multi line string", func() {
			c := compiler.New()

			prog, err := c.Compile(fmt.Sprintf(`
				LET obj = { "foo": "bar", %s: "wsx"}

				RETURN obj["qaz"]
			`, "`qaz`"))

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
