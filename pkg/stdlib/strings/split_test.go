package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestSplit(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Split(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Split(context.Background(), values.NewString("foo"))

			So(err, ShouldBeError)
		})
	})

	Convey("Split('foo-bar-baz', '-' ) should return an array", t, func() {
		out, err := strings.Split(
			context.Background(),
			values.NewString("foo-bar-baz"),
			values.NewString("-"),
		)

		So(err, ShouldBeNil)

		So(out.String(), ShouldEqual, `["foo","bar","baz"]`)
	})

	Convey("Split('foo-bar-baz', '-', 2) should return an array", t, func() {
		out, err := strings.Split(
			context.Background(),
			values.NewString("foo-bar-baz"),
			values.NewString("-"),
			values.NewInt(2),
		)

		So(err, ShouldBeNil)

		So(out.String(), ShouldEqual, `["foo","bar-baz"]`)
	})
}
