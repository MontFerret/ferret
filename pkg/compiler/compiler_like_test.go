package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLikeOperator(t *testing.T) {
	Convey("RETURN \"foo\"  LIKE  \"f*\" ", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			RETURN "foo" LIKE  "f*" 
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `true`)
	})

	Convey("RETURN LIKE('foo', 'f*')", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			RETURN LIKE('foo', 'f*') 
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `true`)
	})

	Convey("RETURN \"foo\" NOT LIKE  \"b*\" ", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			RETURN "foo" NOT LIKE  "b*" 
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `true`)
	})

	Convey("LET t = \"foo\" LIKE  \"f*\" ", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			LET res = "foo" LIKE  "f*"

			RETURN res
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `true`)
	})

	Convey("FOR IN LIKE", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			FOR str IN ["foo", "bar", "qaz"]
				FILTER str LIKE "*a*"
				RETURN str 
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `["bar","qaz"]`)
	})

	Convey("FOR IN LIKE 2", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			FOR str IN ["foo", "bar", "qaz"]
				FILTER str LIKE "*a*"
				RETURN str 
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `["bar","qaz"]`)
	})

	Convey("LIKE ternary", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			RETURN ("foo" NOT LIKE  "b*") ? "foo" : "bar"
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `"foo"`)
	})

	Convey("LIKE ternary 2", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			RETURN true ? ("foo" NOT LIKE  "b*") : false
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `true`)
	})

	Convey("LIKE ternary 3", t, func() {
		c := compiler.New()

		out1, err := c.MustCompile(`
			RETURN true ? false : ("foo" NOT LIKE  "b*")
		`).Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out1), ShouldEqual, `false`)
	})
}
