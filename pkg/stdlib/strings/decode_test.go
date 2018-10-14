package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFromBase64(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.FromBase64(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When hash is not valid base64", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.FromBase64(
				context.Background(),
				values.NewString("foobar"),
			)

			So(err, ShouldBeError)
		})
	})

	Convey("Should decode a given hash", t, func() {
		out, err := strings.FromBase64(
			context.Background(),
			values.NewString("Zm9vYmFy"),
		)

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, "Zm9vYmFy")
		So(out, ShouldEqual, "foobar")
	})
}
