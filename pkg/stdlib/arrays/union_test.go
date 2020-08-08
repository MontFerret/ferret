package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
	. "github.com/smartystreets/goconvey/convey"
)

// TestUnion returns the union of distinct values of all passed arrays.
// @param arrays {Any[], repeated} - List of arrays to combine.
// @return {Any[]} - All array elements combined in a single array, without duplicates, in any order.
func TestUnion(t *testing.T) {
	Convey("Should union all arrays", t, func() {
		arr1 := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
		)

		arr2 := values.NewArrayWith(
			values.NewString("a"),
			values.NewString("b"),
			values.NewString("c"),
			values.NewString("d"),
		)

		arr3 := values.NewArrayWith(
			values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(2),
			),
			values.NewArrayWith(
				values.NewInt(3),
				values.NewInt(4),
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
