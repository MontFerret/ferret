package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestForWhile(t *testing.T) {
	Convey("FOR WHILE", t, func() {
		Convey("WHILE", func() {
			Convey("Should compile FOR i WHILE false RETURN i", func() {
				c := compiler.New()

				p, err := c.Compile(`
			FOR i WHILE false
				RETURN i
		`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, "[]")
			})

			Convey("Should compile FOR i WHILE [] RETURN i", func() {
				c := compiler.New()

				p, err := c.Compile(`
			FOR i WHILE []
				RETURN i
		`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, "[]")
			})

			Convey("Should compile FOR i WHILE F() < 10 RETURN i", func() {
				c := compiler.New()

				counter := -1
				c.RegisterFunction("F", func(ctx context.Context, args ...core.Value) (core.Value, error) {
					counter++
					return values.NewInt(counter), nil
				})

				p, err := c.Compile(`
			FOR i WHILE F() < 10
				RETURN i
		`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, "[0,1,2,3,4,5,6,7,8,9]")
			})

			Convey("Should compile FOR i WHILE F() RETURN i", func() {
				c := compiler.New()

				counter := -1
				c.RegisterFunction("F", func(ctx context.Context, args ...core.Value) (core.Value, error) {
					counter++

					if counter == 10 {
						return values.False, nil
					}

					return values.True, nil
				})

				p, err := c.Compile(`
			FOR i WHILE F()
				RETURN i
		`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, "[0,1,2,3,4,5,6,7,8,9]")
			})

			Convey("Should compile nested FOR operators", func() {
				c := compiler.New()

				counter := -1
				c.RegisterFunction("F", func(ctx context.Context, args ...core.Value) (core.Value, error) {
					counter++
					return values.NewInt(counter), nil
				})

				p, err := c.Compile(`
			FOR i WHILE F() < 5
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`)

				So(err, ShouldBeNil)

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)

				So(string(out), ShouldEqual, "[0,1,2,2,4,6,3,6,9,12,4,8,12,16,20]")
			})
		})

		Convey("DO-WHILE", func() {
			Convey("Should compile FOR i DO WHILE false RETURN i", func() {
				c := compiler.New()

				p, err := c.Compile(`
			FOR i DO WHILE false
				RETURN i
		`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, "[0]")
			})

			Convey("Should compile FOR i DO WHILE [] RETURN i", func() {
				c := compiler.New()

				p, err := c.Compile(`
			FOR i DO WHILE []
				RETURN i
		`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, "[0]")
			})

			Convey("Should compile FOR i DO WHILE F() < 10 RETURN i", func() {
				c := compiler.New()

				counter := -1
				c.RegisterFunction("F", func(ctx context.Context, args ...core.Value) (core.Value, error) {
					counter++
					return values.NewInt(counter), nil
				})

				p, err := c.Compile(`
			FOR i DO WHILE F() < 10
				RETURN i
		`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, "[0,1,2,3,4,5,6,7,8,9,10]")
			})

			Convey("Should compile FOR i WHILE F() RETURN i", func() {
				c := compiler.New()

				counter := -1
				c.RegisterFunction("F", func(ctx context.Context, args ...core.Value) (core.Value, error) {
					counter++

					if counter == 10 {
						return values.False, nil
					}

					return values.True, nil
				})

				p, err := c.Compile(`
			FOR i WHILE F()
				RETURN i
		`)

				So(err, ShouldBeNil)
				So(p, ShouldHaveSameTypeAs, &runtime.Program{})

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, "[0,1,2,3,4,5,6,7,8,9]")
			})

			Convey("Should compile nested FOR operators", func() {
				c := compiler.New()

				counter := -1
				c.RegisterFunction("F", func(ctx context.Context, args ...core.Value) (core.Value, error) {
					counter++
					return values.NewInt(counter), nil
				})

				p, err := c.Compile(`
			FOR i WHILE F() < 5
				LET y = i + 1
				FOR x IN 1..y
					RETURN i * x
		`)

				So(err, ShouldBeNil)

				out, err := p.Run(context.Background())

				So(err, ShouldBeNil)

				So(string(out), ShouldEqual, "[0,1,2,2,4,6,3,6,9,12,4,8,12,16,20]")
			})
		})
	})
}
