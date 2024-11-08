package compiler_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestArrayOperator(t *testing.T) {
	Convey("ALL", t, func() {
		Convey("[1,2,3] ALL IN [1,2,3] should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3] ALL IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[1,2,4] ALL IN [1,2,3] should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,4] ALL IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("[4,5,6] ALL NOT IN [1,2,3] should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [4,5,6] ALL NOT IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[1,2,3] ALL > 0 should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3] ALL > 0
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[1,2,3] ALL > 2 should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3] ALL > 2
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("[1,2,3] ALL >= 3 should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3] ALL >= 3
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("['foo','bar'] ALL != 'moo' should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN ['foo', 'bar'] ALL != 'moo'
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})
	})

	Convey("ANY", t, func() {
		Convey("[1,2,3] ANY IN [1,2,3] should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3] ANY IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[4,2,5] ANY IN [1,2,3] should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [4,2,5] ANY IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[4,5,6] ANY IN [1,2,3] should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [4,5,6] ANY IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("[4,5,6] ANY NOT IN [1,2,3] should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [4,5,6] ANY NOT IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[1,2,3 ] ANY == 2 should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3 ] ANY == 2
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[1,2,3 ] ANY == 4 should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3 ] ANY == 4
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("['foo','bar'] ANY == 'foo' should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN ['foo', 'bar'] ANY == 'foo'
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})
	})

	Convey("NONE", t, func() {
		Convey("[1,2,3] NONE IN [1,2,3] should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3] NONE IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("[4,2,5] NONE IN [1,2,3] should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [4,2,5] NONE IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("[4,5,6] NONE IN [1,2,3] should return true", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [4,5,6] NONE IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[4,5,6] NONE NOT IN [1,2,3] should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [4,5,6] NONE NOT IN [1,2,3]
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("[1,2,3] NONE > 99 should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3] NONE > 99
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `true`)
		})

		Convey("[1,2,3] NONE < 99 should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN [1,2,3] NONE < 99
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})

		Convey("['foo','bar'] NONE == 'foo' should return false", func() {
			c := compiler.New()

			prog, err := c.Compile(`
			RETURN ['foo','bar'] NONE == 'foo'
		`)

			So(err, ShouldBeNil)

			out, err := prog.Run(context.Background())

			So(err, ShouldBeNil)
			So(string(out), ShouldEqual, `false`)
		})
	})
}

func BenchmarkArrayOperatorALL(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [1,2,3] ALL IN [1,2,3]
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorALL2(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [1,2,4] ALL IN [1,2,3]
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorANY(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [1,2,3] ANY IN [1,2,3]
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorANY2(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [4,5,6] ANY IN [1,2,3]
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorANY3(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [4,5,6] ANY NOT IN [1,2,3]
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorANY4(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [1,2,3 ] ANY == 2
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorNONE(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [1,2,3] NONE IN [1,2,3]
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorNONE2(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [4,5,6] NONE IN [1,2,3]
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorNONE3(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [4,5,6] NONE NOT IN [1,2,3]
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}

func BenchmarkArrayOperatorNONE4(b *testing.B) {
	p := compiler.New().MustCompile(`
RETURN [1,2,3] NONE < 99
		`)

	ctx := context.Background()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p.Run(ctx)
	}
}
