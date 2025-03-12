package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnionDistinct(t *testing.T) {
	Convey("Should union all arrays with unique values", t, func() {
		arr1 := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
		)
		arr2 := internal.NewArrayWith(
			core.NewInt(5),
			core.NewInt(2),
			core.NewInt(6),
			core.NewInt(4),
		)

		arr3 := internal.NewArrayWith(
			core.NewString("a"),
			core.NewString("b"),
			core.NewString("c"),
			core.NewString("d"),
		)

		arr4 := internal.NewArrayWith(
			core.NewString("e"),
			core.NewString("b"),
			core.NewString("f"),
			core.NewString("d"),
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
