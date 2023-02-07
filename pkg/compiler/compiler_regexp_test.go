package compiler_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestRegexpOperator(t *testing.T) {
	Convey("Should be possible to use positive regular expression operator", t, func() {
		out := compiler.New().
			MustCompile(`
			RETURN "foo" =~ "^f[o].$" 
		`).
			MustRun(context.Background())

		So(string(out), ShouldEqual, `true`)
	})

	Convey("Should be possible to use negative regular expression operator", t, func() {
		out := compiler.New().
			MustCompile(`
			RETURN "foo" !~ "[a-z]+bar$" 
		`).
			MustRun(context.Background())

		So(string(out), ShouldEqual, `true`)
	})

	Convey("Should be possible to use negative regular expression operator", t, func() {
		c := compiler.New()
		c.Namespace("T").RegisterFunction("REGEXP", func(_ context.Context, _ ...core.Value) (value core.Value, e error) {
			return values.NewString("[a-z]+bar$"), nil
		})

		out := c.
			MustCompile(`
			RETURN "foo" !~ T::REGEXP()
		`).
			MustRun(context.Background())

		So(string(out), ShouldEqual, `true`)
	})

	Convey("Should return an error during compilation when a regexp string invalid", t, func() {
		_, err := compiler.New().
			Compile(`
			RETURN "foo" !~ "[ ]\K(?<!\d )(?=(?: ?\d){8})(?!(?: ?\d){9})\d[ \d]+\d" 
		`)

		So(err, ShouldBeError)
	})

	Convey("Should return an error during compilation when a regexp is not a string", t, func() {
		right := []string{
			"[]",
			"{}",
			"1",
			"1.1",
			"TRUE",
		}

		for _, r := range right {
			_, err := compiler.New().
				Compile(fmt.Sprintf(`
			RETURN "foo" !~ %s 
		`, r))

			So(err, ShouldBeError)
		}
	})
}
