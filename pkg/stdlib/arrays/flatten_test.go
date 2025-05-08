package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestFlatten(t *testing.T) {
	Convey("Should flatten an array with depth 1", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewArrayWith(
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewArrayWith(
					runtime.NewInt(5),
					runtime.NewInt(6),
				),
			),
			runtime.NewInt(7),
			runtime.NewArrayWith(
				runtime.NewInt(8),
				runtime.NewArrayWith(
					runtime.NewInt(9),
					runtime.NewArrayWith(
						runtime.NewInt(10),
					),
				),
			),
		)

		out, err := arrays.Flatten(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,[5,6],7,8,[9,[10]]]")
	})

	Convey("Should flatten an array with depth more than 1", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewArrayWith(
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewArrayWith(
					runtime.NewInt(5),
					runtime.NewInt(6),
				),
			),
			runtime.NewInt(7),
			runtime.NewArrayWith(
				runtime.NewInt(8),
				runtime.NewArrayWith(
					runtime.NewInt(9),
					runtime.NewArrayWith(
						runtime.NewInt(10),
					),
				),
			),
		)

		out, err := arrays.Flatten(context.Background(), arr, runtime.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6,7,8,9,[10]]")
	})
}
