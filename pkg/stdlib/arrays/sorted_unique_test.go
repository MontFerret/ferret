package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestSortedUnique(t *testing.T) {
	Convey("Should sort numbers", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(1),
			core.NewInt(6),
			core.NewInt(2),
			core.NewInt(6),
			core.NewInt(5),
			core.NewInt(1),
			core.NewInt(4),
		)

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})

	Convey("Should sort strings", t, func() {
		arr := internal.NewArrayWith(
			core.NewString("e"),
			core.NewString("b"),
			core.NewString("a"),
			core.NewString("c"),
			core.NewString("a"),
			core.NewString("d"),
			core.NewString("f"),
			core.NewString("d"),
			core.NewString("e"),
			core.NewString("f"),
		)

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["a","b","c","d","e","f"]`)
	})

	Convey("Should return empty array", t, func() {
		arr := internal.NewArrayWith()

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `[]`)
	})
}
