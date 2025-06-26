package types_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/types"
)

func TestToBool(t *testing.T) {
	Convey("Bool", t, func() {
		Convey("When true should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.True,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})

		Convey("When false should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.False,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.False)
		})
	})

	Convey("Int", t, func() {
		Convey("When > 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewInt(1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})

		Convey("When < 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewInt(-1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})

		Convey("When 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.ZeroInt,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.False)
		})
	})

	Convey("Float", t, func() {
		Convey("When > 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewFloat(1.1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})

		Convey("When < 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewFloat(-1.1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})

		Convey("When 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.ZeroFloat,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.False)
		})
	})

	Convey("String", t, func() {
		Convey("When != '' should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewString("foobar"),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})

		Convey("When == '' should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewString(""),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.False)
		})
	})

	Convey("DateTime", t, func() {
		Convey("When > 0  should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewCurrentDateTime(),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})

		Convey("When == 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.ZeroDateTime,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.False)
		})
	})

	Convey("None", t, func() {
		Convey("Should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.None,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.False)
		})
	})

	Convey("arrayList", t, func() {
		Convey("Should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewArray(0),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})
	})

	Convey("hashMap", t, func() {
		Convey("Should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				runtime.NewObject(),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.True)
		})
	})
}
