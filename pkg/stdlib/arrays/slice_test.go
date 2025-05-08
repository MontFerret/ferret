package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestSlice(t *testing.T) {
	Convey("Should return a sliced array with a given start position ", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, runtime.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[4,5,6]")
	})

	Convey("Should return an empty array when start position is out of bounds", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, runtime.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})

	Convey("Should return a sliced array with a given start position and length", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.Slice(
			context.Background(),
			arr,
			runtime.NewInt(2),
			runtime.NewInt(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,4]")
	})

	Convey("Should return an empty array when length is out of bounds", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
			runtime.NewInt(6),
		)

		out, err := arrays.Slice(context.Background(), arr, runtime.NewInt(2), runtime.NewInt(10))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[3,4,5,6]")
	})
}
