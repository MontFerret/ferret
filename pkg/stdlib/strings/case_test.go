package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestLower(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Lower(context.Background())

			So(err, ShouldBeError)

		})
	})

	Convey("Lower('FOOBAR') should return 'foobar'", t, func() {
		out, _ := strings.Lower(
			context.Background(),
			values.NewString("FOOBAR"),
		)

		So(out, ShouldEqual, "foobar")
	})
}

func TestUpper(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Upper(context.Background())

			So(err, ShouldBeError)

		})
	})

	Convey("Lower('foobar') should return 'FOOBAR'", t, func() {
		out, _ := strings.Upper(
			context.Background(),
			values.NewString("foobar"),
		)

		So(out, ShouldEqual, "FOOBAR")
	})
}
