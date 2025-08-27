package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestIntersection_Basic(t *testing.T) {
	ctx := context.Background()

	Convey("Should find intersections between 2 arrays", t, func() {
		arr1 := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		arr2 := runtime.NewArrayWith(
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
			runtime.NewInt(7),
			runtime.NewInt(8),
			runtime.NewInt(9),
		)

		out, err := arrays.Intersection(ctx, arr1, arr2)

		check := map[int]bool{
			4: true,
			5: true,
			6: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*runtime.Array)

		size, _ := arr.Length(ctx)

		So(size, ShouldEqual, 3)

		_ = arr.ForEach(ctx, func(_ context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			_, exists := check[int(value.(runtime.Int))]

			So(exists, ShouldBeTrue)

			return true, nil
		})
	})

	Convey("Should find intersections between more than 2 arrays", t, func() {
		arr1 := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		arr2 := runtime.NewArrayWith(
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		arr3 := runtime.NewArrayWith(
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
			runtime.NewInt(7),
		)

		out, err := arrays.Intersection(ctx, arr1, arr2, arr3)

		check := map[int]bool{
			3: true,
			4: true,
			5: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*runtime.Array)
		size, _ := arr.Length(ctx)

		So(size, ShouldEqual, 3)

		_ = arr.ForEach(ctx, func(_ context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			_, exists := check[int(value.(runtime.Int))]

			So(exists, ShouldBeTrue)

			return true, nil
		})
	})
}

func TestIntersection_EdgeCases(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle empty arrays", t, func() {
		arr1 := runtime.NewArrayWith()
		arr2 := runtime.NewArrayWith(runtime.NewInt(1))
		out, err := arrays.Intersection(ctx, arr1, arr2)
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should handle identical arrays", t, func() {
		arr1 := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
		arr2 := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
		out, err := arrays.Intersection(ctx, arr1, arr2)
		So(err, ShouldBeNil)
		// Should contain both elements
		length, lengthErr := out.(runtime.List).Length(ctx)
		So(lengthErr, ShouldBeNil)
		So(length, ShouldEqual, 2)
	})
}

func TestIntersection_ArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Should reject too few arguments", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(1))
		_, err := arrays.Intersection(ctx, arr)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject invalid argument types", t, func() {
		nonArray := runtime.NewString("not an array")
		arr := runtime.NewArrayWith(runtime.NewInt(1))

		_, err := arrays.Intersection(ctx, nonArray, arr)
		So(err, ShouldNotBeNil)
	})
}
