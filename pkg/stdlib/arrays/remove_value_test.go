package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveValue_Basic(t *testing.T) {
	ctx := context.Background()

	Convey("Should return a copy of an array without given element(s)", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(3),
		)

		out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4]")
	})

	Convey("Should return a copy of an array without given element(s) with limit", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(3),
			runtime.NewInt(5),
			runtime.NewInt(3),
		)

		out, err := arrays.RemoveValue(
			ctx,
			arr,
			runtime.NewInt(3),
			runtime.Int(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4,5,3]")
	})
}

func TestRemoveValue_EdgeCases(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle limit of 0", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(0))
		So(err, ShouldBeNil)
		// With limit 0, nothing should be removed
		So(out.String(), ShouldEqual, "[1,2,2,3]")
	})

	Convey("Should handle negative limit", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(-1))
		So(err, ShouldBeNil)
		// Negative limit should remove all occurrences
		So(out.String(), ShouldEqual, "[1,3]")

		out, err = arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(-10))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,3]")
	})

	Convey("Should handle removing non-existent value", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(999))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3]")
	})

	Convey("Should handle positive limit correctly", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)
		out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(1))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3]")
	})
}

func TestRemoveValue_ArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Should reject too few arguments", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(1))
		_, err := arrays.RemoveValue(ctx, arr)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject invalid argument types", t, func() {
		nonArray := runtime.NewString("not an array")
		_, err := arrays.RemoveValue(ctx, nonArray, runtime.NewInt(1))
		So(err, ShouldNotBeNil)
	})
}

func TestRemoveValue_SpecialValues(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle None values correctly", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.None,
			runtime.NewInt(3),
		)

		out, err := arrays.RemoveValue(ctx, arr, runtime.None)
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,3]")
	})
}
