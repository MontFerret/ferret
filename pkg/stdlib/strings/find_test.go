package strings_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestFindFirst(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.FindFirst(context.Background())

			So(err, ShouldBeError)

			_, err = strings.FindFirst(
				context.Background(),
				core.NewString("foo"),
			)

			So(err, ShouldBeError)
		})
	})

	Convey("When args are strings", t, func() {
		Convey("FindFirst('foobarbaz', 'ba') should return 3", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				core.NewString("foobarbaz"),
				core.NewString("ba"),
			)

			So(out, ShouldEqual, 3)
		})

		Convey("FindFirst('foobarbaz', 'ba', 4) should return 6", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				core.NewString("foobarbaz"),
				core.NewString("ba"),
				core.NewInt(4),
			)

			So(out, ShouldEqual, 6)
		})

		Convey("FindFirst('foobarbaz', 'ba', 4) should return -1", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				core.NewString("foobarbaz"),
				core.NewString("ba"),
				core.NewInt(7),
			)

			So(out, ShouldEqual, -1)
		})

		Convey("FindFirst('foobarbaz', 'ba', 0, 3) should return -1", func() {
			out, _ := strings.FindFirst(
				context.Background(),
				core.NewString("foobarbaz"),
				core.NewString("ba"),
				core.NewInt(0),
				core.NewInt(3),
			)

			So(out, ShouldEqual, -1)
		})
	})
}

func TestFindLast(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.FindLast(context.Background())

			So(err, ShouldBeError)

			_, err = strings.FindLast(
				context.Background(),
				core.NewString("foo"),
			)

			So(err, ShouldBeError)
		})
	})

	Convey("When args are strings", t, func() {
		Convey("FindLast('foobarbaz', 'ba') should return 6", func() {
			out, _ := strings.FindLast(
				context.Background(),
				core.NewString("foobarbaz"),
				core.NewString("ba"),
			)

			So(out, ShouldEqual, 6)
		})

		Convey("FindLast('foobarbaz', 'ba', 7) should return -1", func() {
			out, _ := strings.FindLast(
				context.Background(),
				core.NewString("foobarbaz"),
				core.NewString("ba"),
				core.NewInt(7),
			)

			So(out, ShouldEqual, -1)
		})

		Convey("FindLast('foobarbaz', 'ba', 0, 5) should return 3", func() {
			out, _ := strings.FindLast(
				context.Background(),
				core.NewString("foobarbaz"),
				core.NewString("ba"),
				core.NewInt(0),
				core.NewInt(5),
			)

			So(out, ShouldEqual, 3)
		})
	})
}
