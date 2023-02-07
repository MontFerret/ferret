package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestContains(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := strings.Contains(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When args are strings", t, func() {
		Convey("Contains('foobar', 'foo') should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				values.NewString("foobar"),
				values.NewString("bar"),
			)

			So(out, ShouldEqual, values.True)
		})

		Convey("Contains('foobar', 'qaz') should return 'false'", func() {
			out, _ := strings.Contains(
				context.Background(),
				values.NewString("foobar"),
				values.NewString("qaz"),
			)

			So(out, ShouldEqual, values.False)
		})

		Convey("Contains('foobar', 'foo', true) should return '3'", func() {
			out, _ := strings.Contains(
				context.Background(),
				values.NewString("foobar"),
				values.NewString("bar"),
				values.True,
			)

			So(out, ShouldEqual, 3)
		})

		Convey("Contains('foobar', 'qaz', true) should return '-1'", func() {
			out, _ := strings.Contains(
				context.Background(),
				values.NewString("foobar"),
				values.NewString("bar"),
				values.True,
			)

			So(out, ShouldEqual, 3)
		})
	})

	Convey("When args are not strings", t, func() {
		Convey("Contains('foo123', 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				values.NewString("foo123"),
				values.NewInt(1),
			)

			So(out, ShouldEqual, values.True)
		})

		Convey("Contains(123, 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				values.NewInt(123),
				values.NewInt(1),
			)

			So(out, ShouldEqual, values.True)
		})

		Convey("Contains([1,2,3], 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				values.NewArrayWith(values.NewInt(1), values.NewInt(2), values.NewInt(3)),
				values.NewInt(1),
			)

			So(out, ShouldEqual, values.True)
		})
	})
}
