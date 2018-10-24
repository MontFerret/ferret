package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFunctionCall(t *testing.T) {
	Convey("Should compile RETURN TYPENAME(1)", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN TYPENAME(1)
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `"int"`)
	})

	Convey("Should compile WAIT(10) RETURN 1", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			WAIT(10)
			RETURN 1
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `1`)
	})

	Convey("Should compile LET duration = 10 WAIT(duration) RETURN 1", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			LET duration = 10

			WAIT(duration)

			RETURN 1
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `1`)
	})

	Convey("Should compile function call inside FOR IN statement", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			FOR i IN [1, 2, 3, 4]
				LET duration = 10

				WAIT(duration)

				RETURN i * 2
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, `[2,4,6,8]`)
	})
}

func BenchmarkFunctionCallArg1(b *testing.B) {
	c := compiler.New()

	c.RegisterFunction("TEST", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return values.None, nil
	})

	p := c.MustCompile(`
			RETURN TYPENAME(1)
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkFunctionCallArg2(b *testing.B) {
	c := compiler.New()

	c.RegisterFunction("TEST", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return values.None, nil
	})

	p := c.MustCompile(`
			RETURN TYPENAME(1, 2)
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkFunctionCallArg3(b *testing.B) {
	c := compiler.New()

	c.RegisterFunction("TEST", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return values.None, nil
	})

	p := c.MustCompile(`
			RETURN TYPENAME(1, 2, 3)
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkFunctionEmpty(b *testing.B) {
	c := compiler.New()
	c.RegisterFunction("TEST", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return values.None, nil
	})

	p := c.MustCompile(`
			RETURN TEST()
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
