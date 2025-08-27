package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestNth_Basic(t *testing.T) {
	ctx := context.Background()

	Convey("Should return item by index", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Nth(ctx, arr, runtime.NewInt(1))

		So(err, ShouldBeNil)
		So(out.(runtime.Comparable).Compare(runtime.NewInt(2)), ShouldEqual, 0)
	})

	Convey("Should return None when no value", t, func() {
		arr := runtime.NewArrayWith()

		out, err := arrays.Nth(ctx, arr, runtime.NewInt(1))

		So(err, ShouldBeNil)
		So(out, ShouldPointTo, runtime.None)
	})

	Convey("Should return None when passed negative value", t, func() {
		arr := runtime.NewArrayWith()

		out, err := arrays.Nth(ctx, arr, runtime.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldPointTo, runtime.None)
	})
}

func TestNth_EdgeCases(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle very large index", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(10),
			runtime.NewInt(20),
			runtime.NewInt(30),
		)
		out, err := arrays.Nth(ctx, arr, runtime.NewInt(1000))
		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})

	Convey("Should handle negative index correctly", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(10),
			runtime.NewInt(20),
			runtime.NewInt(30),
		)
		out, err := arrays.Nth(ctx, arr, runtime.NewInt(-1))
		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)

		out, err = arrays.Nth(ctx, arr, runtime.NewInt(-100))
		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})

	Convey("Should handle valid indices correctly", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(10),
			runtime.NewInt(20),
			runtime.NewInt(30),
		)
		out, err := arrays.Nth(ctx, arr, runtime.NewInt(0))
		So(err, ShouldBeNil)
		So(out.(runtime.Comparable).Compare(runtime.NewInt(10)), ShouldEqual, 0)

		out, err = arrays.Nth(ctx, arr, runtime.NewInt(2))
		So(err, ShouldBeNil)
		So(out.(runtime.Comparable).Compare(runtime.NewInt(30)), ShouldEqual, 0)
	})

	Convey("Should handle index at boundary", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(10),
			runtime.NewInt(20),
			runtime.NewInt(30),
		)
		out, err := arrays.Nth(ctx, arr, runtime.NewInt(3))
		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}

func TestNth_ArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Should reject too few arguments", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(1))
		_, err := arrays.Nth(ctx, arr)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject invalid argument types", t, func() {
		nonArray := runtime.NewString("not an array")
		nonInt := runtime.NewString("not an int")
		arr := runtime.NewArrayWith(runtime.NewInt(1))

		_, err := arrays.Nth(ctx, nonArray, runtime.NewInt(0))
		So(err, ShouldNotBeNil)

		_, err = arrays.Nth(ctx, arr, nonInt)
		So(err, ShouldNotBeNil)
	})
}

func TestNth_SpecialValues(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle None values correctly", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.None,
			runtime.NewInt(3),
		)

		noneVal, err := arrays.Nth(ctx, arr, runtime.NewInt(1))
		So(err, ShouldBeNil)
		So(noneVal, ShouldEqual, runtime.None)
	})
}
