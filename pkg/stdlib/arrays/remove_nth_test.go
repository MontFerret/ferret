package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveNth(t *testing.T) {
	Convey("Should return a copy of an array without an element by its position", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.RemoveNth(context.Background(), arr, values.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4,5]")
	})

	Convey("Should return a copy of an array with all elements when a position is invalid", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.RemoveNth(context.Background(), arr, values.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5]")
	})
}
