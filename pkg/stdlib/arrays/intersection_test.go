package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestIntersection(t *testing.T) {
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

		out, err := arrays.Intersection(context.Background(), arr1, arr2)

		check := map[int]bool{
			4: true,
			5: true,
			6: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*runtime.Array)

		size, _ := arr.Length(context.Background())

		So(size, ShouldEqual, 3)

		_ = arr.ForEach(context.Background(), func(_ context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
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

		out, err := arrays.Intersection(context.Background(), arr1, arr2, arr3)

		check := map[int]bool{
			3: true,
			4: true,
			5: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*runtime.Array)
		size, _ := arr.Length(context.Background())

		So(size, ShouldEqual, 3)

		_ = arr.ForEach(context.Background(), func(_ context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			_, exists := check[int(value.(runtime.Int))]

			So(exists, ShouldBeTrue)

			return true, nil
		})
	})
}
