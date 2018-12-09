package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPrecedence(t *testing.T) {
	Convey("Math operators", t, func() {
		Convey("2 + 2 * 2", func() {
			c := compiler.New()

			p := c.MustCompile(`RETURN 2 + 2 * 2`)

			out := p.MustRun(context.Background())

			So(string(out), ShouldEqual, "6")
		})

		Convey("2 * 2 + 2", func() {
			c := compiler.New()

			p := c.MustCompile(`RETURN 2 * 2 + 2`)

			out := p.MustRun(context.Background())

			So(string(out), ShouldEqual, "6")
		})
	})
}
