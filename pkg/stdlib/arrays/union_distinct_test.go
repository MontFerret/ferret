package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnionDistinct(t *testing.T) {
	Convey("Should union all arrays with unique values", t, func() {
		arr1 := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
		)
		arr2 := values.NewArrayWith(
			values.NewInt(5),
			values.NewInt(2),
			values.NewInt(6),
			values.NewInt(4),
		)

		arr3 := values.NewArrayWith(
			values.NewString("a"),
			values.NewString("b"),
			values.NewString("c"),
			values.NewString("d"),
		)

		arr4 := values.NewArrayWith(
			values.NewString("e"),
			values.NewString("b"),
			values.NewString("f"),
			values.NewString("d"),
		)

		out, err := arrays.UnionDistinct(
			context.Background(),
			arr1,
			arr2,
			arr3,
			arr4,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `[1,2,3,4,5,6,"a","b","c","d","e","f"]`)
	})
}
