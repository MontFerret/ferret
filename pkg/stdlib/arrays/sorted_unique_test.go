package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestSortedUnique(t *testing.T) {
	Convey("Should sort numbers", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(1),
			runtime.NewInt(6),
			runtime.NewInt(2),
			runtime.NewInt(6),
			runtime.NewInt(5),
			runtime.NewInt(1),
			runtime.NewInt(4),
		)

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})

	Convey("Should sort strings", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewString("e"),
			runtime.NewString("b"),
			runtime.NewString("a"),
			runtime.NewString("c"),
			runtime.NewString("a"),
			runtime.NewString("d"),
			runtime.NewString("f"),
			runtime.NewString("d"),
			runtime.NewString("e"),
			runtime.NewString("f"),
		)

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["a","b","c","d","e","f"]`)
	})

	Convey("Should return empty array", t, func() {
		arr := runtime.NewArrayWith()

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `[]`)
	})
}
