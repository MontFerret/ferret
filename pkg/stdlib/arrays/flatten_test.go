package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestFlatten_Basic(t *testing.T) {
	ctx := context.Background()

	Convey("Should flatten an array with depth 1", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewArrayWith(
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewArrayWith(
					runtime.NewInt(5),
					runtime.NewInt(6),
				),
			),
			runtime.NewInt(7),
			runtime.NewArrayWith(
				runtime.NewInt(8),
				runtime.NewArrayWith(
					runtime.NewInt(9),
					runtime.NewArrayWith(
						runtime.NewInt(10),
					),
				),
			),
		)

		out, err := arrays.Flatten(ctx, arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,[5,6],7,8,[9,[10]]]")
	})

	Convey("Should flatten an array with depth more than 1", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewArrayWith(
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewArrayWith(
					runtime.NewInt(5),
					runtime.NewInt(6),
				),
			),
			runtime.NewInt(7),
			runtime.NewArrayWith(
				runtime.NewInt(8),
				runtime.NewArrayWith(
					runtime.NewInt(9),
					runtime.NewArrayWith(
						runtime.NewInt(10),
					),
				),
			),
		)

		out, err := arrays.Flatten(ctx, arr, runtime.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6,7,8,9,[10]]")
	})
}

func TestFlatten_EdgeCases(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle empty arrays", t, func() {
		arr := runtime.NewArrayWith()
		out, err := arrays.Flatten(ctx, arr)
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should handle depth 0 correctly", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewArrayWith(runtime.NewInt(1)),
			runtime.NewArrayWith(runtime.NewInt(2)),
		)
		out, err := arrays.Flatten(ctx, arr, runtime.NewInt(0))
		So(err, ShouldBeNil)
		// With depth 0, no flattening should occur
		So(out.String(), ShouldEqual, "[[1],[2]]")
	})

	Convey("Should handle negative depth", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewArrayWith(runtime.NewInt(1)),
			runtime.NewArrayWith(runtime.NewInt(2)),
		)
		out, err := arrays.Flatten(ctx, arr, runtime.NewInt(-1))
		So(err, ShouldBeNil)
		// Negative depth should not flatten
		So(out.String(), ShouldEqual, "[[1],[2]]")
	})

	Convey("Should handle very deep nesting", t, func() {
		// Create [[[[[5]]]]]
		deep := runtime.Value(runtime.NewInt(5))
		for i := 0; i < 5; i++ {
			deep = runtime.NewArrayWith(deep)
		}
		arr := runtime.NewArrayWith(deep)

		out, err := arrays.Flatten(ctx, arr, runtime.NewInt(6))
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[5]")
	})

	Convey("Should handle mixed types correctly", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewString("hello"),
			runtime.NewArrayWith(runtime.NewString("world")),
			runtime.NewBoolean(true),
			runtime.NewArrayWith(runtime.None),
		)
		out, err := arrays.Flatten(ctx, arr)
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["hello","world",true,null]`)
	})
}

func TestFlatten_ArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Should reject invalid argument types", t, func() {
		nonArray := runtime.NewString("not an array")
		_, err := arrays.Flatten(ctx, nonArray)
		So(err, ShouldNotBeNil)
	})
}
