package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestIntersection(t *testing.T) {
	Convey("Should find intersections between 2 arrays", t, func() {
		arr1 := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		arr2 := values.NewArrayWith(
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
			values.NewInt(7),
			values.NewInt(8),
			values.NewInt(9),
		)

		out, err := arrays.Intersection(context.Background(), arr1, arr2)

		check := map[int]bool{
			4: true,
			5: true,
			6: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*values.Array)

		So(arr.Length(), ShouldEqual, 3)

		arr.ForEach(func(value core.Value, idx int) bool {
			_, exists := check[int(value.(values.Int))]

			So(exists, ShouldBeTrue)

			return true
		})
	})

	Convey("Should find intersections between more than 2 arrays", t, func() {
		arr1 := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		arr2 := values.NewArrayWith(
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		arr3 := values.NewArrayWith(
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
			values.NewInt(7),
		)

		out, err := arrays.Intersection(context.Background(), arr1, arr2, arr3)

		check := map[int]bool{
			3: true,
			4: true,
			5: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*values.Array)

		So(arr.Length(), ShouldEqual, 3)

		arr.ForEach(func(value core.Value, idx int) bool {
			_, exists := check[int(value.(values.Int))]

			So(exists, ShouldBeTrue)

			return true
		})
	})
}
