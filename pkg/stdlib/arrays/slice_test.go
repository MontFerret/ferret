package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestSlice_Basic(t *testing.T) {
	ctx := context.Background()

	Convey("Should return a sliced array with a given start position ", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.Slice(ctx, arr, runtime.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[4,5,6]")
	})

	Convey("Should return an empty array when start position is out of bounds", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.Slice(ctx, arr, runtime.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should return a sliced array with a given start position and length", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.Slice(
			ctx,
			arr,
			runtime.NewInt(2),
			runtime.NewInt(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,4]")
	})

	Convey("Should return an empty array when length is out of bounds", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.Slice(ctx, arr, runtime.NewInt(2), runtime.NewInt(10))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,4,5,6]")
	})
}

func TestSlice_EdgeCases(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle negative start index", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.Slice(ctx, arr, runtime.NewInt(-1))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")

		out, err = arrays.Slice(ctx, arr, runtime.NewInt(-10))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should handle negative length", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.Slice(ctx, arr, runtime.NewInt(1), runtime.NewInt(-1))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")

		out, err = arrays.Slice(ctx, arr, runtime.NewInt(0), runtime.NewInt(-5))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should handle zero length", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.Slice(ctx, arr, runtime.NewInt(1), runtime.NewInt(0))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should handle start index beyond array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.Slice(ctx, arr, runtime.NewInt(10))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should handle length beyond array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.Slice(ctx, arr, runtime.NewInt(1), runtime.NewInt(100))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[2,3]")
	})
}

func TestSlice_ArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Should reject too few arguments", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(1))
		_, err := arrays.Slice(ctx, arr)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject invalid argument types", t, func() {
		nonArray := runtime.NewString("not an array")
		nonInt := runtime.NewString("not an int")
		arr := runtime.NewArrayWith(runtime.NewInt(1))

		_, err := arrays.Slice(ctx, nonArray, runtime.NewInt(0))
		So(err, ShouldNotBeNil)

		_, err = arrays.Slice(ctx, arr, nonInt)
		So(err, ShouldNotBeNil)
	})
}
