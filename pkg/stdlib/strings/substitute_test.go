package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestSubstitute(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Substitute(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Substitute(context.Background(), values.NewString("foo"))

			So(err, ShouldBeError)
		})
	})

	Convey("Substitute('foo-bar-baz', 'a', 'o') should return 'foo-bor-boz'", t, func() {
		out, err := strings.Substitute(
			context.Background(),
			values.NewString("foo-bar-baz"),
			values.NewString("a"),
			values.NewString("o"),
		)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, "foo-bor-boz")
	})

	Convey("Substitute('foo-bar-baz', 'a', 'o', 1) should return 'foo-bor-baz'", t, func() {
		out, err := strings.Substitute(
			context.Background(),
			values.NewString("foo-bar-baz"),
			values.NewString("a"),
			values.NewString("o"),
			values.NewInt(1),
		)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, "foo-bor-baz")
	})
}
