package compiler_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/pkg/stdlib/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUseExpression(t *testing.T) {
	c := compiler.New()

	c.Namespace("X").
		RegisterFunctions(core.NewFunctionsFromMap(
			map[string]core.Function{
				"XXX_CONTAINS": strings.Contains,
				"XXX_UPPER":    strings.Upper,
			},
		))

	c.Namespace("Z").
		RegisterFunction("XXX_TO_STRING", types.ToString)

	// Y contain the same function as Z to test for future collistions
	c.Namespace("Y").
		RegisterFunction("XXX_TO_STRING", types.ToString)

	Convey("Use Expression", t, func() {

		Convey("Should compile", func() {

			Convey("Single statement", func() {
				p, err := c.Compile(`
				USE X
	
				RETURN XXX_CONTAINS("s", "s")
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "true")
			})

			Convey("Many functions from one lib", func() {
				p, err := c.Compile(`
				USE X
	
				RETURN XXX_CONTAINS(XXX_UPPER("s"), "S")
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "true")
			})

			Convey("Many statements", func() {
				p, err := c.Compile(`
				USE X
				USE Z
	
				RETURN XXX_TO_STRING(XXX_CONTAINS("s", "s"))
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, `"true"`)
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
					// Z and Y both contain function "XXX_CONTAINS"
					USE Z
					USE Y

					RETURN 1`,
				},
				{
					Name: "USE namespace twice",
					Query: `
					USE X
					USE X

					RETURN XXX_CONTAINS("s", "s")
					`,
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
