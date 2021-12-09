package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestForLimit(t *testing.T) {
	Convey("Should compile query with LIMIT 2", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				LIMIT 2
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[1,2]`)
	})

	Convey("Should compile query with LIMIT 2, 2", t, func() {
		c := compiler.New()

		// 4 is offset
		// 2 is count
		p, err := c.Compile(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT 4, 2
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[5,6]`)
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
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LET x = "foo"
				TEST(x)
				LIMIT 2
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(counter, ShouldEqual, 2)
		So(string(out), ShouldEqual, `[1,2]`)
	})

	Convey("Should be able to reuse values from a source", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LET x = i
				LIMIT 2
				RETURN i*x
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[1,4]`)
	})

	Convey("Should be able to use variable", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET li = 2
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT li
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[1,2]`)
	})

	Convey("Should be able to use function call", t, func() {
		c := compiler.New()
		c.RegisterFunction("TEST", func(ctx context.Context, args ...core.Value) (core.Value, error) {
			return values.NewInt(2), nil
		})

		p, err := c.Compile(`
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT TEST()
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[1,2]`)
	})

	Convey("Should be able to use member expression (object)", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET o = {
				limit: 2
			}
			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT o.limit
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[1,2]`)
	})

	Convey("Should be able to use member expression (array)", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET o = [1,2]

			FOR i IN [ 1,2,3,4,5,6,7,8 ]
				LIMIT o[1]
				RETURN i
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[1,2]`)
	})
}
