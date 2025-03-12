package strings_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestLike(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Like(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Like(context.Background(), core.NewString(""))

			So(err, ShouldBeError)
		})
	})

	Convey("Should return true when matches with _ pattern", t, func() {
		out, _ := strings.Like(
			context.Background(),
			core.NewString("cart"),
			core.NewString("ca_t"),
		)

		So(out, ShouldEqual, true)
	})

	Convey("Should return true when matches with % pattern", t, func() {
		out, _ := strings.Like(
			context.Background(),
			core.NewString("foo bar baz"),
			core.NewString("%bar%"),
		)

		So(out, ShouldEqual, true)
	})

	Convey("Should return false when matches with no caseInsensitive parameter", t, func() {
		out, _ := strings.Like(
			context.Background(),
			core.NewString("FoO bAr BaZ"),
			core.NewString("fOo%bAz"),
		)

		So(out, ShouldEqual, false)
	})

	Convey("Should return true when matches with caseInsensitive parameter", t, func() {
		out, _ := strings.Like(
			context.Background(),
			core.NewString("FoO bAr BaZ"),
			core.NewString("fOo%bAz"),
			core.True,
		)

		So(out, ShouldEqual, true)
	})
}
