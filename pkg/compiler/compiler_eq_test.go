package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEqualityOperators(t *testing.T) {
	Convey("Should compile RETURN 2 > 1", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			RETURN 2 > 1
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, "true")
	})
}
