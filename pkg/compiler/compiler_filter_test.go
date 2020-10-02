package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestForFilter(t *testing.T) {
	Convey("Should compile query with FILTER i > 2", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 2
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[3,4,3]`)
	})

	Convey("Should compile query with FILTER i > 1 AND i < 3", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 1 AND i < 4
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,3,3]`)
	})

	Convey("Should compile query with a regexp FILTER statement", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET users = [
				{
					age: 31,
					gender: "m",
					name: "Josh"
				},
				{
					age: 29,
					gender: "f",
					name: "Mary"
				},
				{
					age: 36,
					gender: "m",
					name: "Peter"
				}
			]
			FOR u IN users
				FILTER u.name =~ "r"
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"age":29,"gender":"f","name":"Mary"},{"age":36,"gender":"m","name":"Peter"}]`)
	})

	Convey("Should compile query with multiple FILTER statements", t, func() {
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
				FILTER u.active == true
				FILTER u.age < 35
				RETURN u
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m"},{"active":true,"age":29,"gender":"f"}]`)
	})

	Convey("Should compile query with multiple FILTER statements", t, func() {
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

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m"}]`)
	})

	Convey("Should compile query with left side expression", t, func() {
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

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m"},{"active":true,"age":29,"gender":"f"},{"active":true,"age":36,"gender":"m"}]`)
	})

	Convey("Should compile query with multiple left side expression", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
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

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":true,"age":31,"gender":"m","married":true},{"active":true,"age":45,"gender":"f","married":true}]`)
	})

	Convey("Should compile query with multiple left side expression and with binary operator", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
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

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[{"active":false,"age":69,"gender":"m","married":true}]`)
	})

	Convey("Should compile query with multiple left side expression and with binary operator 2", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
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

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[]`)
	})

	Convey("Should define variables", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LET x = 2
				FILTER i > x
				RETURN i + x
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[5,6,5]`)
	})

	Convey("Should call functions", t, func() {
		c := compiler.New()
		counterA := 0
		counterB := 0
		c.RegisterFunction("COUNT_A", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counterA++

			return values.None, nil
		})

		c.RegisterFunction("COUNT_B", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			counterB++

			return values.None, nil
		})

		p, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LET x = 2
				COUNT_A()
				FILTER i > x
				COUNT_B()
				RETURN i + x
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(counterA, ShouldEqual, 6)
		So(counterB, ShouldEqual, 3)
		So(string(out), ShouldEqual, `[5,6,5]`)
	})
}

func BenchmarkFilter(b *testing.B) {
	p := compiler.New().MustCompile(`
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
				FILTER u.age < 35
				RETURN u
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkFilter2(b *testing.B) {
	p := compiler.New().MustCompile(`
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
				FILTER u.age < 35
				RETURN u
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkFilter3(b *testing.B) {
	p := compiler.New().MustCompile(`
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

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkFilter4(b *testing.B) {
	p := compiler.New().MustCompile(`
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
				LIMIT 2
				RETURN u
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
