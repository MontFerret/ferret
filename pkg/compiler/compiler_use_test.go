package compiler_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/pkg/stdlib/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUseExpression(t *testing.T) {
	c := compiler.New()

	c.Namespace("X").
		RegisterFunction("INCLUDES", strings.Contains)
	c.Namespace("Z").
		RegisterFunction("STRING_FROM", types.ToString)
	c.Namespace("Y").
		RegisterFunction("INCLUDES", strings.Contains)

	Convey("Use Expression", t, func() {

		Convey("Should compile", func() {

			Convey("Single statement", func() {
				p, err := c.Compile(`
				USE X
	
				RETURN INCLUDES("s", "s")
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "true")
			})

			Convey("Many statements", func() {
				p, err := c.Compile(`
				USE X
				USE Z
	
				RETURN STRING_FROM(INCLUDES("s", "s"))
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, `"true"`)
			})

			Convey("USE namespace twice", func() {
				p, err := c.Compile(`
				USE X
				USE X
	
				RETURN INCLUDES("s", "s")
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "true")
			})

			Convey("Namespace doesn't exists", func() {
				p, err := c.Compile(`
				USE NOT::EXISTS
		
				RETURN 1
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "1")
			})
		})

		Convey("Should not compile", func() {

			testCases := []struct {
				Name  string
				Query string
			}{
				{
					Name: "Wrong namespace format",
					Query: `
					USE NOT::EXISTS::
			
					RETURN 1`,
				},
				{
					Name: "Empty namespace",
					Query: `
					USE
			
					RETURN 1`,
				},
				{
					Name: "Functions collision",
					Query: `
					// X and Y both contain function "INCLUDES"
					USE X
					USE Y

					RETURN 1`,
				},
			}

			for _, tC := range testCases {
				Convey(tC.Name, func() {
					_, err := c.Compile(tC.Query)
					So(err, ShouldNotBeNil)
				})
			}
		})
	})
}
