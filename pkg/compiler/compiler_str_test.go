package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	Convey("Should be possible to use multi line string", t, func() {
		out := compiler.New().
			MustCompile(`
			RETURN "
FOO
BAR
"
		`).
			MustRun(context.Background())

		So(string(out), ShouldEqual, `"\nFOO\nBAR\n"`)
	})
}
