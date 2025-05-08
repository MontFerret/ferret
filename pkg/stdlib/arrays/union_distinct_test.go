package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnionDistinct(t *testing.T) {
	Convey("Should union all arrays with unique values", t, func() {
		arr1 := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
		)
		arr2 := runtime.NewArrayWith(
			runtime.NewInt(5),
			runtime.NewInt(2),
			runtime.NewInt(6),
			runtime.NewInt(4),
		)

		arr3 := runtime.NewArrayWith(
			runtime.NewString("a"),
			runtime.NewString("b"),
			runtime.NewString("c"),
			runtime.NewString("d"),
		)

		arr4 := runtime.NewArrayWith(
			runtime.NewString("e"),
			runtime.NewString("b"),
			runtime.NewString("f"),
			runtime.NewString("d"),
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
