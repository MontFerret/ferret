package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestShift(t *testing.T) {
	Convey("Should return a copy of an array without the first element", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Shift(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[2,3,4,5]")
	})

	Convey("Should return empty array if a given one is empty", t, func() {
		arr := internal.NewArray(0)

		out, err := arrays.Shift(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})
}
