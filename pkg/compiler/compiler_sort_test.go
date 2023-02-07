package compiler_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestForSort(t *testing.T) {
	Convey("Should compile query with SORT statement", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
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
				SORT u.age
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":29,"gender":"f"},{"active":true,"age":31,"gender":"m"},{"active":true,"age":36,"gender":"m"}]`)
	})

	Convey("Should compile query with SORT DESC statement", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
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
				SORT u.age DESC
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":36,"gender":"m"},{"active":true,"age":31,"gender":"m"},{"active":true,"age":29,"gender":"f"}]`)
	})

	Convey("Should compile query with SORT statement with multiple expressions", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
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
					age: 31,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				SORT u.age, u.gender
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":29,"gender":"f"},{"active":true,"age":31,"gender":"f"},{"active":true,"age":31,"gender":"m"},{"active":true,"age":36,"gender":"m"}]`)
	})

	Convey("Should define variables and call functions", t, func() {
		c := compiler.New()
		counter := 0
		c.RegisterFunction("TEST", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counter++

			So(args[0], ShouldEqual, "foo")

			return values.None, nil
		})

		p, err := c.Compile(`
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
					age: 31,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				LET x = "foo"
				TEST(x)
				SORT u.age, u.gender
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(counter, ShouldEqual, 4)
		So(string(out), ShouldEqual, `[{"active":true,"age":29,"gender":"f"},{"active":true,"age":31,"gender":"f"},{"active":true,"age":31,"gender":"m"},{"active":true,"age":36,"gender":"m"}]`)
	})

	Convey("Should be able to reuse values from a source", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
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
					age: 31,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				}
			]
			FOR u IN users
				LET x = u.gender
				SORT u.age, u.gender
				RETURN {u,x}
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[{"u":{"active":true,"age":29,"gender":"f"},"x":"f"},{"u":{"active":true,"age":31,"gender":"f"},"x":"f"},{"u":{"active":true,"age":31,"gender":"m"},"x":"m"},{"u":{"active":true,"age":36,"gender":"m"},"x":"m"}]`)
	})
}
