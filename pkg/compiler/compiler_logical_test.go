package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLogicalOperators(t *testing.T) {
	Convey("Should compile RETURN 2 > 1 AND 1 > 0", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 2 > 1 AND 1 > 0
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("Should compile RETURN 2 > 1 OR 1 < 0", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 2 > 1 OR 1 < 0
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})

	Convey("1 || 7  should return 1", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 1 || 7
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "1")
	})

	Convey("NONE || 'foo'  should return 'foo'", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN NONE || 'foo'
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `"foo"`)
	})

	Convey("NONE && true  should return null", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN NONE && true
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `null`)
	})

	Convey("'' && true  should return ''", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN '' && true
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `""`)
	})

	Convey("true && 23  should return '23", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN true && 23 
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `23`)
	})

	Convey("NOT TRUE should return false", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN NOT TRUE
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `false`)
	})

	Convey("NOT u.valid should return true", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET u = { valid: false }

			RETURN NOT u.valid
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `true`)
	})
}
