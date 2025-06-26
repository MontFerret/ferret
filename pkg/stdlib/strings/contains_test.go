package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"

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
				core.NewString("foobar"),
				core.NewString("bar"),
			)

			So(out, ShouldEqual, core.True)
		})

		Convey("Contains('foobar', 'qaz') should return 'false'", func() {
			out, _ := strings.Contains(
				context.Background(),
				core.NewString("foobar"),
				core.NewString("qaz"),
			)

			So(out, ShouldEqual, core.False)
		})

		Convey("Contains('foobar', 'foo', true) should return '3'", func() {
			out, _ := strings.Contains(
				context.Background(),
				core.NewString("foobar"),
				core.NewString("bar"),
				core.True,
			)

			So(out, ShouldEqual, 3)
		})

		Convey("Contains('foobar', 'qaz', true) should return '-1'", func() {
			out, _ := strings.Contains(
				context.Background(),
				core.NewString("foobar"),
				core.NewString("bar"),
				core.True,
			)

			So(out, ShouldEqual, 3)
		})
	})

	Convey("When args are not strings", t, func() {
		Convey("Contains('foo123', 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				core.NewString("foo123"),
				core.NewInt(1),
			)

			So(out, ShouldEqual, core.True)
		})

		Convey("Contains(123, 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				core.NewInt(123),
				core.NewInt(1),
			)

			So(out, ShouldEqual, core.True)
		})

		Convey("Contains([1,2,3], 1) should return 'true'", func() {
			out, _ := strings.Contains(
				context.Background(),
				runtime.NewArrayWith(core.NewInt(1), core.NewInt(2), core.NewInt(3)),
				core.NewInt(1),
			)

			So(out, ShouldEqual, core.True)
		})
	})
}
