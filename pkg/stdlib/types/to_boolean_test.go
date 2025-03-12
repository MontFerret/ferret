package types_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/types"
)

func TestToBool(t *testing.T) {
	Convey("Bool", t, func() {
		Convey("When true should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				core.True,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})

		Convey("When false should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				core.False,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.False)
		})
	})

	Convey("Int", t, func() {
		Convey("When > 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				core.NewInt(1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})

		Convey("When < 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				core.NewInt(-1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})

		Convey("When 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				core.ZeroInt,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.False)
		})
	})

	Convey("Float", t, func() {
		Convey("When > 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				core.NewFloat(1.1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})

		Convey("When < 0 should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				core.NewFloat(-1.1),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})

		Convey("When 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				core.ZeroFloat,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.False)
		})
	})

	Convey("String", t, func() {
		Convey("When != '' should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				core.NewString("foobar"),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})

		Convey("When == '' should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				core.NewString(""),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.False)
		})
	})

	Convey("DateTime", t, func() {
		Convey("When > 0  should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				core.NewCurrentDateTime(),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})

		Convey("When == 0 should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				core.ZeroDateTime,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.False)
		})
	})

	Convey("None", t, func() {
		Convey("Should return false", func() {
			out, err := types.ToBool(
				context.Background(),
				core.None,
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.False)
		})
	})

	Convey("Array", t, func() {
		Convey("Should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				internal.NewArray(0),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})
	})

	Convey("Object", t, func() {
		Convey("Should return true", func() {
			out, err := types.ToBool(
				context.Background(),
				internal.NewObject(),
			)

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})
	})
}
