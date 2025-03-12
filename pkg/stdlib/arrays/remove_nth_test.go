package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveNth(t *testing.T) {
	Convey("Should return a copy of an array without an element by its position", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.RemoveNth(context.Background(), arr, core.NewInt(2))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4,5]")
	})

	Convey("Should return a copy of an array with all elements when a position is invalid", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.RemoveNth(context.Background(), arr, core.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5]")
	})
}
