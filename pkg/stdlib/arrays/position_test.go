package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPosition(t *testing.T) {
	Convey("Should return TRUE when a value exists in a given array", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Position(context.Background(), arr, core.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "true")
	})

	Convey("Should return FALSE when a value does not exist in a given array", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Position(context.Background(), arr, core.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "false")
	})

	Convey("Should return index when a value exists in a given array", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Position(
			context.Background(),
			arr,
			core.NewInt(3),
			core.NewBoolean(true),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "2")
	})

	Convey("Should return -1 when a value does not exist in a given array", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Position(
			context.Background(),
			arr,
			core.NewInt(6),
			core.NewBoolean(true),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "-1")
	})
}
