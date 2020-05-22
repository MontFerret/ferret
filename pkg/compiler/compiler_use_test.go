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

	newCompiler := func() *compiler.Compiler {
		c := compiler.New()

		err := c.Namespace("X").
			RegisterFunctions(core.NewFunctionsFromMap(
				map[string]core.Function{
					"XXX_CONTAINS": strings.Contains,
					"XXX_UPPER":    strings.Upper,
				},
			))
		So(err, ShouldBeNil)

		return c
	}

	Convey("Use Expression", t, func() {

		Convey("Should compile", func() {

			Convey("Single statement", func() {
				p, err := newCompiler().Compile(`
				USE X
	
				RETURN XXX_CONTAINS("s", "s")
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "true")
			})

			Convey("Single statement (lower case)", func() {
				p, err := newCompiler().Compile(`
				use x
	
				return xxx_contains("s", "s")
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "true")
			})

			Convey("Many functions from one lib", func() {
				p, err := newCompiler().Compile(`
				USE X
	
				RETURN XXX_CONTAINS(XXX_UPPER("s"), "S")
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "true")
			})

			Convey("Many statements", func() {
				c := newCompiler()

				// Z must contain functions different from X
				c.Namespace("Z").
					RegisterFunction("XXX_TO_STRING", types.ToString)

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
				p, err := newCompiler().Compile(`
				USE NOT::EXISTS
		
				RETURN 1
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "1")
			})

			Convey("Full and short path works together", func() {
				p, err := newCompiler().Compile(`
				USE X

				LET short = XXX_CONTAINS("s", "s")
				LET full = X::XXX_CONTAINS("s", "s")

				RETURN short == full
				`)

				So(err, ShouldBeNil)
				out := p.MustRun(context.Background())

				So(string(out), ShouldEqual, "true")
			})

			Convey("Nested namespace", func() {
				c := newCompiler()

				c.Namespace("X").
					Namespace("Y").
					RegisterFunction("YYY_CONTAINS", strings.Contains)

				Convey("Short path", func() {
					p, err := c.Compile(`
					USE X
		
					RETURN Y::YYY_CONTAINS("s", "s")
					`)

					So(err, ShouldBeNil)
					out := p.MustRun(context.Background())

					So(string(out), ShouldEqual, "true")
				})

				Convey("Full path", func() {
					p, err := c.Compile(`
					USE X
		
					RETURN X::Y::YYY_CONTAINS("s", "s")
					`)

					So(err, ShouldBeNil)
					out := p.MustRun(context.Background())

					So(string(out), ShouldEqual, "true")
				})
			})
		})

		Convey("Should not compile", func() {

			testCases := []struct {
				Name     string
				Query    string
				Compiler func() *compiler.Compiler
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
					Name: "USE namespace twice",
					Query: `
					USE X
					USE X

					RETURN XXX_CONTAINS("s", "s")
					`,
				},
				{
					Name: "Functions collision",
					Query: `
					// Z and Y both contain function "XXX_CONTAINS"
					USE Z
					USE Y

					RETURN 1`,
					Compiler: func() *compiler.Compiler {
						c := newCompiler()

						c.Namespace("Z").
							RegisterFunction("XXX_TO_STRING", types.ToString)

						// Y contain the same function as Z to test for future collistions
						c.Namespace("Y").
							RegisterFunction("XXX_TO_STRING", types.ToString)

						return c
					},
				},
				{
					Name: "Nested namespace collision",
					Query: `
					USE X
					USE Z

					RETURN 1
					`,
					Compiler: func() *compiler.Compiler {
						c := newCompiler()

						c.Namespace("X").
							Namespace("Y").
							RegisterFunction("YYY_CONTAINS", strings.Contains)

						c.Namespace("Z").
							Namespace("Y").
							RegisterFunction("YYY_CONTAINS", strings.Contains)

						return c
					},
				},
			}

			for _, tC := range testCases {
				Convey(tC.Name, func() {
					c := newCompiler()
					if tC.Compiler != nil {
						c = tC.Compiler()
					}

					_, err := c.Compile(tC.Query)
					So(err, ShouldNotBeNil)
				})
			}
		})
	})
}
