package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnshift(t *testing.T) {
	Convey("Should return a copy of an array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Unshift(context.Background(), arr, runtime.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")
	})

	Convey("Should ignore non-unique items", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Unshift(
			context.Background(),
			arr,
			runtime.NewInt(0),
			runtime.True,
		)

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")

		out2, err := arrays.Unshift(
			context.Background(),
			arr,
			runtime.NewInt(0),
			runtime.True,
		)

		So(err, ShouldBeNil)
		So(out2, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")
	})
}
