package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCollect(t *testing.T) {
	Convey("Should not have access to initial variables", t, func() {
		c := compiler.New()

		_, err := c.Compile(`
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
			FOR i IN users
				COLLECT gender = i.gender
				RETURN {
					user: i,
					gender: gender
				}
		`)

		So(err, ShouldNotBeNil)
	})

	Convey("Should group result by a single key", t, func() {
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
			FOR i IN users
				COLLECT gender = i.gender
				RETURN gender
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `["f","m"]`)
	})

	Convey("Should group result by multiple keys", t, func() {
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
			FOR i IN users
				COLLECT gender = i.gender, age = i.age
				RETURN {gender, age}
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `["f","m"]`)
	})
}