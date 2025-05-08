package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPosition(t *testing.T) {
	Convey("Should return TRUE when a value exists in a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Position(context.Background(), arr, runtime.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "true")
	})

	Convey("Should return FALSE when a value does not exist in a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Position(context.Background(), arr, runtime.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "false")
	})

	Convey("Should return index when a value exists in a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Position(
			context.Background(),
			arr,
			runtime.NewInt(3),
			runtime.NewBoolean(true),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "2")
	})

	Convey("Should return -1 when a value does not exist in a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Position(
			context.Background(),
			arr,
			runtime.NewInt(6),
			runtime.NewBoolean(true),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "-1")
	})
}
