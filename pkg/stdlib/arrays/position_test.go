package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPosition_Basic(t *testing.T) {
	ctx := context.Background()

	Convey("Should return TRUE when a value exists in a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Position(ctx, arr, runtime.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "true")
	})

	Convey("Should return FALSE when a value does not exist in a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Position(ctx, arr, runtime.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "false")
	})

	Convey("Should return index when a value exists in a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Position(
			ctx,
			arr,
			runtime.NewInt(3),
			runtime.NewBoolean(true),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "2")
	})

	Convey("Should return -1 when a value does not exist in a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Position(
			ctx,
			arr,
			runtime.NewInt(6),
			runtime.NewBoolean(true),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "-1")
	})
}

func TestPosition_EdgeCases(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle searching for None", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.None,
			runtime.NewInt(3),
		)
		out, err := arrays.Position(ctx, arr, runtime.None)
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "true")
	})

	Convey("Should handle empty array", t, func() {
		arr := runtime.NewArrayWith()
		out, err := arrays.Position(ctx, arr, runtime.NewInt(1))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "false")
	})

	Convey("Should return correct index with position flag", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.Position(ctx, arr, runtime.NewInt(2), runtime.NewBoolean(true))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "1")
	})

	Convey("Should return correct index for None", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.None,
			runtime.NewInt(3),
		)
		out, err := arrays.Position(ctx, arr, runtime.None, runtime.NewBoolean(true))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "1")
	})
}

func TestPosition_ArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Should reject too few arguments", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(1))
		_, err := arrays.Position(ctx, arr)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject invalid argument types", t, func() {
		nonArray := runtime.NewString("not an array")
		nonBool := runtime.NewString("not a bool")
		arr := runtime.NewArrayWith(runtime.NewInt(1))

		_, err := arrays.Position(ctx, nonArray, runtime.NewInt(1))
		So(err, ShouldNotBeNil)

		_, err = arrays.Position(ctx, arr, runtime.NewInt(1), nonBool)
		So(err, ShouldNotBeNil)
	})
}
