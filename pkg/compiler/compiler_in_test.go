package compiler_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestInOperator(t *testing.T) {
	Convey("1 IN [1,2,3] should return true", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 1 IN [1,2,3]
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `true`)
	})

	Convey("4 IN [1,2,3] should return false", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 4 IN [1,2,3]
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `false`)
	})

	Convey("4 NOT IN [1,2,3] should return true", t, func() {
		c := compiler.New()

		p, err := c.Compile(`
			RETURN 4 NOT IN [1,2,3]
		`)

		So(err, ShouldBeNil)

		out, err := p.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `true`)
	})
}

func BenchmarkInOperator(b *testing.B) {
	p := compiler.New().MustCompile(`
			RETURN 1 IN [1,2,3]
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkInOperatorNot(b *testing.B) {
	p := compiler.New().MustCompile(`
			RETURN 4 NOT IN [1,2,3]
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
