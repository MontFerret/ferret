package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestOutersection(t *testing.T) {
	Convey("Should find intersections between 2 arrays", t, func() {
		arr1 := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)

		arr2 := runtime.NewArrayWith(
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
		)

		out, err := arrays.Outersection(context.Background(), arr1, arr2)

		check := map[int]bool{
			1: true,
			4: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*runtime.Array)
		size, _ := arr.Length(context.Background())

		So(size, ShouldEqual, 2)

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
		)

		arr2 := runtime.NewArrayWith(
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
		)

		arr3 := runtime.NewArrayWith(
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Outersection(context.Background(), arr1, arr2, arr3)

		check := map[int]bool{
			1: true,
			5: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*runtime.Array)
		size, _ := arr.Length(context.Background())

		So(size, ShouldEqual, 2)

		_ = arr.ForEach(context.Background(), func(_ context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			_, exists := check[int(value.(runtime.Int))]

			So(exists, ShouldBeTrue)

			return true, nil
		})
	})
}
