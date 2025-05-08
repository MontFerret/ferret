package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

// TestUnion returns the union of distinct values of all passed arrays.
// @param arrays {Any[], repeated} - List of arrays to combine.
// @return {Any[]} - All array elements combined in a single array, without duplicates, in any order.
func TestUnion(t *testing.T) {
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

		out, err := arrays.Union(
			context.Background(),
			arr1,
			arr2,
			arr3,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `[1,2,3,4,"a","b","c","d",[1,2],[3,4]]`)
	})
}
