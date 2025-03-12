package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestIntersection(t *testing.T) {
	Convey("Should find intersections between 2 arrays", t, func() {
		arr1 := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
		)

		arr2 := internal.NewArrayWith(
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
			core.NewInt(7),
			core.NewInt(8),
			core.NewInt(9),
		)

		out, err := arrays.Intersection(context.Background(), arr1, arr2)

		check := map[int]bool{
			4: true,
			5: true,
			6: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*internal.Array)

		So(arr.Length(), ShouldEqual, 3)

		arr.ForEach(func(value core.Value, idx int) bool {
			_, exists := check[int(value.(core.Int))]

			So(exists, ShouldBeTrue)

			return true
		})
	})

	Convey("Should find intersections between more than 2 arrays", t, func() {
		arr1 := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		arr2 := internal.NewArrayWith(
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
		)

		arr3 := internal.NewArrayWith(
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
			core.NewInt(7),
		)

		out, err := arrays.Intersection(context.Background(), arr1, arr2, arr3)

		check := map[int]bool{
			3: true,
			4: true,
			5: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*internal.Array)

		So(arr.Length(), ShouldEqual, 3)

		arr.ForEach(func(value core.Value, idx int) bool {
			_, exists := check[int(value.(core.Int))]

			So(exists, ShouldBeTrue)

			return true
		})
	})
}
