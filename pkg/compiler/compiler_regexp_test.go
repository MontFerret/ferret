package compiler_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
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
}
