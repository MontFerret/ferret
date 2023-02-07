package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestSlice(t *testing.T) {
	Convey("Should return a sliced array with a given start position ", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, values.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[4,5,6]")
	})

	Convey("Should return an empty array when start position is out of bounds", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, values.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should return a sliced array with a given start position and length", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		out, err := arrays.Slice(
			context.Background(),
			arr,
			values.NewInt(2),
			values.NewInt(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,4]")
	})

	Convey("Should return an empty array when length is out of bounds", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, values.NewInt(2), values.NewInt(10))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,4,5,6]")
	})
}
