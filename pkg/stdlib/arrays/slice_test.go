package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestSlice(t *testing.T) {
	Convey("Should return a sliced array with a given start position ", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, core.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[4,5,6]")
	})

	Convey("Should return an empty array when start position is out of bounds", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, core.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should return a sliced array with a given start position and length", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
		)

		out, err := arrays.Slice(
			context.Background(),
			arr,
			core.NewInt(2),
			core.NewInt(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,4]")
	})

	Convey("Should return an empty array when length is out of bounds", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
			core.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, core.NewInt(2), core.NewInt(10))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,4,5,6]")
	})
}
