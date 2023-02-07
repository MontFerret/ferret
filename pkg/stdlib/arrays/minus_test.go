package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestMinus(t *testing.T) {
	Convey("Should find differences between 2 arrays", t, func() {
		arr1 := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
		)

		arr2 := values.NewArrayWith(
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		out, err := arrays.Minus(context.Background(), arr1, arr2)

		check := map[int]bool{
			1: true,
			2: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*values.Array)

		So(arr.Length(), ShouldEqual, 2)

		arr.ForEach(func(value core.Value, idx int) bool {
			_, exists := check[int(value.(values.Int))]

			So(exists, ShouldBeTrue)

			return true
		})
	})

	Convey("Should find differences between more than 2 arrays", t, func() {
		arr1 := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
		)

		arr2 := values.NewArrayWith(
			values.NewInt(3),
			values.NewInt(9),
			values.NewInt(5),
			values.NewInt(6),
		)

		arr3 := values.NewArrayWith(
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
			values.NewInt(7),
			values.NewInt(8),
		)

		out, err := arrays.Minus(context.Background(), arr1, arr2, arr3)

		check := map[int]bool{
			1: true,
			2: true,
		}

		So(err, ShouldBeNil)

		arr := out.(*values.Array)

		So(arr.Length(), ShouldEqual, 2)

		arr.ForEach(func(value core.Value, idx int) bool {
			_, exists := check[int(value.(values.Int))]

			So(exists, ShouldBeTrue)

			return true
		})
	})
}
