package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

// TestConcat_Basic returns the union of distinct values of all passed arrays.
// @param arrays {Any[], repeated} List of arrays to combine.
// @return {Any[]} All array elements combined in a single array, without duplicates, in any order.
func TestConcat_Basic(t *testing.T) {
	ctx := context.Background()

	Convey("Should union all arrays", t, func() {
		arr1 := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
		)

		arr2 := runtime.NewArrayWith(
			runtime.NewString("a"),
			runtime.NewString("b"),
			runtime.NewString("c"),
			runtime.NewString("d"),
		)

		arr3 := runtime.NewArrayWith(
			runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
			),
			runtime.NewArrayWith(
				runtime.NewInt(3),
				runtime.NewInt(4),
			),
		)

		out, err := arrays.Concat(
			ctx,
			arr1,
			arr2,
			arr3,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `[1,2,3,4,"a","b","c","d",[1,2],[3,4]]`)
	})
}

func TestConcat_EdgeCases(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle empty arrays", t, func() {
		arr1 := runtime.NewArrayWith()
		arr2 := runtime.NewArrayWith()
		out, err := arrays.Concat(ctx, arr1, arr2)
		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})
}

func TestConcat_ArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Should reject too few arguments", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(1))
		_, err := arrays.Concat(ctx, arr)
		So(err, ShouldNotBeNil)
	})

	Convey("Should reject invalid argument types", t, func() {
		nonArray := runtime.NewString("not an array")
		arr := runtime.NewArrayWith(runtime.NewInt(1))

		_, err := arrays.Concat(ctx, nonArray, arr)
		So(err, ShouldNotBeNil)
	})
}
