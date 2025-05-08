package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPush(t *testing.T) {
	Convey("Should create a new array with a new element in the end", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Push(context.Background(), arr, runtime.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})

	Convey("Should not add a new element if not unique when uniqueness check is enabled", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Push(
			context.Background(),
			arr,
			runtime.NewInt(6),
			runtime.True,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")

		out2, err := arrays.Push(
			context.Background(),
			arr,
			runtime.NewInt(6),
			runtime.True,
		)

		So(err, ShouldBeNil)
		So(out2.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})
}
