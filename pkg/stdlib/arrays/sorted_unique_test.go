package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestSortedUnique(t *testing.T) {
	Convey("Should sort numbers", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(1),
			values.NewInt(6),
			values.NewInt(2),
			values.NewInt(6),
			values.NewInt(5),
			values.NewInt(1),
			values.NewInt(4),
		)

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})

	Convey("Should sort strings", t, func() {
		arr := values.NewArrayWith(
			values.NewString("e"),
			values.NewString("b"),
			values.NewString("a"),
			values.NewString("c"),
			values.NewString("a"),
			values.NewString("d"),
			values.NewString("f"),
			values.NewString("d"),
			values.NewString("e"),
			values.NewString("f"),
		)

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["a","b","c","d","e","f"]`)
	})

	Convey("Should return empty array", t, func() {
		arr := values.NewArrayWith()

		out, err := arrays.SortedUnique(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `[]`)
	})
}
