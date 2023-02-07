package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPosition(t *testing.T) {
	Convey("Should return TRUE when a value exists in a given array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Position(context.Background(), arr, values.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "true")
	})

	Convey("Should return FALSE when a value does not exist in a given array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Position(context.Background(), arr, values.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "false")
	})

	Convey("Should return index when a value exists in a given array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Position(
			context.Background(),
			arr,
			values.NewInt(3),
			values.NewBoolean(true),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "2")
	})

	Convey("Should return -1 when a value does not exist in a given array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Position(
			context.Background(),
			arr,
			values.NewInt(6),
			values.NewBoolean(true),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "-1")
	})
}
