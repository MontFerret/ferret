package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestFlatten(t *testing.T) {
	Convey("Should flatten an array with depth 1", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewArrayWith(
				values.NewInt(3),
				values.NewInt(4),
				values.NewArrayWith(
					values.NewInt(5),
					values.NewInt(6),
				),
			),
			values.NewInt(7),
			values.NewArrayWith(
				values.NewInt(8),
				values.NewArrayWith(
					values.NewInt(9),
					values.NewArrayWith(
						values.NewInt(10),
					),
				),
			),
		)

		out, err := arrays.Flatten(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,[5,6],7,8,[9,[10]]]")
	})

	Convey("Should flatten an array with depth more than 1", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewArrayWith(
				values.NewInt(3),
				values.NewInt(4),
				values.NewArrayWith(
					values.NewInt(5),
					values.NewInt(6),
				),
			),
			values.NewInt(7),
			values.NewArrayWith(
				values.NewInt(8),
				values.NewArrayWith(
					values.NewInt(9),
					values.NewArrayWith(
						values.NewInt(10),
					),
				),
			),
		)

		out, err := arrays.Flatten(context.Background(), arr, values.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6,7,8,9,[10]]")
	})
}
