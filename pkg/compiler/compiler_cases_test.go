package compiler_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCases(t *testing.T) {

	Convey("Should compile", t, func() {

		Convey("String variable should be case-sensitive", func() {
			param := "String"
			query := fmt.Sprintf(`
				let a = "%s"
				return a
			`, param)

			res := compiler.New().
				MustCompile(query).
				MustRun(context.Background())

			So(string(res), ShouldEqual, `"String"`)
		})

		Convey("Should work at any case", func() {

			testCases := []struct {
				Query  string
				Result string
			}{
				// Key-word only
				{
					Query:  "return 2 * 2",
					Result: "4",
				},
				// With function
				{
					Query:  "return last([7, 2, 99])",
					Result: "99",
				},
				// With namespace
				{
					Query:  "return x::last([7, 2, 99])",
					Result: "99",
				},
			}

			strCases := map[string]func(string) string{
				"Mixed": func(s string) string {
					// Capitalize string.
					// Source: https://stackoverflow.com/questions/33696034/make-first-letter-of-words-uppercase-in-a-string
					// return strings.Title(strings.ToLower(s))
					return cases.Title(language.English).String(strings.ToLower(s))
				},
				"Upper": strings.ToUpper,
				"Lower": strings.ToLower,
			}

			for _, tC := range testCases {
				tC := tC

				for strCase, toCase := range strCases {
					query := toCase(tC.Query)
					tcName := fmt.Sprintf(`%s: "%s"`, strCase, query)

					Convey(tcName, func() {

						compiler := compiler.New()
						compiler.Namespace("X").RegisterFunction("LAST", arrays.Last)

						res := compiler.
							MustCompile(query).
							MustRun(context.Background())

						So(string(res), ShouldEqual, tC.Result)
					})
				}
			}
		})
	})

	Convey("Should not compile", t, func() {

		Convey("Variable should be case-sensitive", func() {
			_, err := compiler.New().
				Compile(
					`
					let a = 1
					return A
				`)

			So(err, ShouldNotBeNil)
		})
	})
}
