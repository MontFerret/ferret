package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCollectInto(t *testing.T) {
	Convey("Should create default projection", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender INTO genders
				RETURN {
					gender,
					values: genders
				}
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[{"gender":"f","values":[{"i":{"active":true,"age":25,"gender":"f","married":false}},{"i":{"active":true,"age":45,"gender":"f","married":true}}]},{"gender":"m","values":[{"i":{"active":true,"age":31,"gender":"m","married":true}},{"i":{"active":true,"age":36,"gender":"m","married":false}},{"i":{"active":false,"age":69,"gender":"m","married":true}}]}]`)
	})

	Convey("Should create custom projection", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender INTO genders = { active: i.active }
				RETURN {
					gender,
					values: genders
				}
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[{"gender":"f","values":[{"active":true},{"active":true}]},{"gender":"m","values":[{"active":true},{"active":true},{"active":false}]}]`)
	})

	Convey("Should create custom projection grouped by multiple keys", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender, age = i.age INTO genders = { active: i.active }
				RETURN {
					age,
					gender,
					values: genders
				}
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[{"age":25,"gender":"f","values":[{"active":true}]},{"age":45,"gender":"f","values":[{"active":true}]},{"age":31,"gender":"m","values":[{"active":true}]},{"age":36,"gender":"m","values":[{"active":true}]},{"age":69,"gender":"m","values":[{"active":false}]}]`)
	})
}

func BenchmarkCollectInto(b *testing.B) {
	p := compiler.New().MustCompile(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender INTO genders
				RETURN {
					gender,
					values: genders
				}
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkCollectInto2(b *testing.B) {
	p := compiler.New().MustCompile(`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender INTO genders = { active: i.active }
				RETURN {
					gender,
					values: genders
				}
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
