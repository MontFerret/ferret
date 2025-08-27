package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestLike(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Like(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Like(context.Background(), runtime.NewString(""))

			So(err, ShouldBeError)
		})
	})

	Convey("Should return true when matches with _ pattern", t, func() {
		out, _ := strings.Like(
			context.Background(),
			runtime.NewString("cart"),
			runtime.NewString("ca_t"),
		)

		So(out, ShouldEqual, runtime.True)
	})

	Convey("Should return true when matches with % pattern", t, func() {
		out, _ := strings.Like(
			context.Background(),
			runtime.NewString("foo bar baz"),
			runtime.NewString("%bar%"),
		)

		So(out, ShouldEqual, runtime.True)
	})

	Convey("Should return false when matches with no caseInsensitive parameter", t, func() {
		out, _ := strings.Like(
			context.Background(),
			runtime.NewString("FoO bAr BaZ"),
			runtime.NewString("fOo%bAz"),
		)

		So(out, ShouldEqual, runtime.False)
	})

	Convey("Should return true when matches with caseInsensitive parameter", t, func() {
		out, _ := strings.Like(
			context.Background(),
			runtime.NewString("FoO bAr BaZ"),
			runtime.NewString("fOo%bAz"),
			runtime.True,
		)

		So(out, ShouldEqual, runtime.True)
	})
}
