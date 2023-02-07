package types_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/types"
)

func TestToBool(t *testing.T) {
	Convey("Bool", t, func() {
		Convey("When true should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.True,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("When false should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				values.False,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.False)
		})
	})

	Convey("Int", t, func() {
		Convey("When > 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewInt(1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("When < 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewInt(-1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("When 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				values.ZeroInt,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.False)
		})
	})

	Convey("Float", t, func() {
		Convey("When > 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewFloat(1.1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("When < 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewFloat(-1.1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("When 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				values.ZeroFloat,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.False)
		})
	})

	Convey("String", t, func() {
		Convey("When != '' should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewString("foobar"),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("When == '' should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewString(""),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.False)
		})
	})

	Convey("DateTime", t, func() {
		Convey("When > 0  should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewCurrentDateTime(),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})

		Convey("When == 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				values.ZeroDateTime,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.False)
		})
	})

	Convey("None", t, func() {
		Convey("Should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				values.None,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.False)
		})
	})

	Convey("Array", t, func() {
		Convey("Should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewArray(0),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})
	})

	Convey("Object", t, func() {
		Convey("Should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				values.NewObject(),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})
	})
}
