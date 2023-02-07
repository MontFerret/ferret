package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestRandomToken(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RandomToken(context.Background())

			So(err, ShouldBeError)

		})
	})

	Convey("When args are invalid", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RandomToken(context.Background(), values.NewString("foo"))

			So(err, ShouldBeError)

		})
	})

	Convey("Should generate random string", t, func() {
		str1, _ := strings.RandomToken(
			context.Background(),
			values.NewInt(8),
		)

		So(str1, ShouldHaveLength, 8)

		str2, _ := strings.RandomToken(
			context.Background(),
			values.NewInt(8),
		)

		So(str2, ShouldHaveLength, 8)

		So(str1, ShouldNotEqual, str2)
	})
}
