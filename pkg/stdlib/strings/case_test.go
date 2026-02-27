package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func TestLower(t *testing.T) {
	Convey("Lower('FOOBAR') should return 'foobar'", t, func() {
		out, _ := strings.Lower(
			context.Background(),
			runtime.NewString("FOOBAR"),
		)

		So(out.String(), ShouldEqual, "foobar")
	})
}

func TestUpper(t *testing.T) {
	Convey("Lower('foobar') should return 'FOOBAR'", t, func() {
		out, _ := strings.Upper(
			context.Background(),
			runtime.NewString("foobar"),
		)

		So(out.String(), ShouldEqual, "FOOBAR")
	})
}
