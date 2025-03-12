package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestFlatten(t *testing.T) {
	Convey("Should flatten an array with depth 1", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			internal.NewArrayWith(
				core.NewInt(3),
				core.NewInt(4),
				internal.NewArrayWith(
					core.NewInt(5),
					core.NewInt(6),
				),
			),
			core.NewInt(7),
			internal.NewArrayWith(
				core.NewInt(8),
				internal.NewArrayWith(
					core.NewInt(9),
					internal.NewArrayWith(
						core.NewInt(10),
					),
				),
			),
		)

		out, err := arrays.Flatten(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,[5,6],7,8,[9,[10]]]")
	})

	Convey("Should flatten an array with depth more than 1", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			internal.NewArrayWith(
				core.NewInt(3),
				core.NewInt(4),
				internal.NewArrayWith(
					core.NewInt(5),
					core.NewInt(6),
				),
			),
			core.NewInt(7),
			internal.NewArrayWith(
				core.NewInt(8),
				internal.NewArrayWith(
					core.NewInt(9),
					internal.NewArrayWith(
						core.NewInt(10),
					),
				),
			),
		)

		out, err := arrays.Flatten(context.Background(), arr, core.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6,7,8,9,[10]]")
	})
}
