package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestForFilter(t *testing.T) {
	Convey("Should compile query with FILTER i > 2", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 2
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[3,4,3]`)
	})

	Convey("Should compile query with FILTER i > 1 AND i < 3", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 1 AND i < 4
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,3,3]`)
	})

	Convey("Should compile query with multiple FILTER statements", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active == true
				FILTER u.age < 35
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m"},{"active":true,"age":29,"gender":"f"}]`)
	})

	Convey("Should compile query with multiple FILTER statements", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active == true
				LIMIT 2
				FILTER u.gender == "m"
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m"}]`)
	})

	Convey("Should compile query with left side expression", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m"},{"active":true,"age":29,"gender":"f"},{"active":true,"age":36,"gender":"m"}]`)
	})

	Convey("Should compile query with multiple left side expression", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER u.active AND u.married
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m","married":true},{"active":true,"age":45,"gender":"f","married":true}]`)
	})

	Convey("Should compile query with multiple left side expression and with binary operator", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER !u.active AND u.married
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":false,"age":69,"gender":"m","married":true}]`)
	})

	Convey("Should compile query with multiple left side expression and with binary operator 2", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
				FILTER !u.active AND !u.married
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[]`)
	})
}
