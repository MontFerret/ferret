package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

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
				runtime.NewString("foobar"),
				runtime.NewString("bar"),
			)

			So(out, ShouldEqual, runtime.True)
		})

		Convey("Contains('foobar', 'qaz') should return 'false'", func() {
			out, _ := strings.Contains(
				context.Background(),
				runtime.NewString("foobar"),
				runtime.NewString("qaz"),
			)

			So(out, ShouldEqual, runtime.False)
		})

		Convey("Contains('foobar', 'foo', true) should return '3'", func() {
			out, _ := strings.Contains(
				context.Background(),
				runtime.NewString("foobar"),
				runtime.NewString("bar"),
				runtime.True,
			)

			So(out, ShouldEqual, 3)
		})

		Convey("Contains('foobar', 'qaz', true) should return '-1'", func() {
			out, _ := strings.Contains(
				context.Background(),
				runtime.NewString("foobar"),
				runtime.NewString("qaz"),
				runtime.True,
			)

			So(out, ShouldEqual, -1)
		})
	})

	Convey("When args are not strings", t, func() {
		Convey("Contains('foo123', 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				runtime.NewString("foo123"),
				runtime.NewInt(1),
			)

			So(out, ShouldEqual, runtime.True)
		})

		Convey("Contains(123, 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				runtime.NewInt(123),
				runtime.NewInt(1),
			)

			So(out, ShouldEqual, runtime.True)
		})

		Convey("Contains([1,2,3], 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
				runtime.NewInt(1),
			)

			So(out, ShouldEqual, runtime.True)
		})
	})
}
