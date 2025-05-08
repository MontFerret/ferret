package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveValue(t *testing.T) {
	Convey("Should return a copy of an array without given element(s)", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(3),
		)

		out, err := arrays.RemoveValue(context.Background(), arr, runtime.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4]")
	})

	Convey("Should return a copy of an array without given element(s) with limit", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(3),
			runtime.NewInt(5),
			runtime.NewInt(3),
		)

		out, err := arrays.RemoveValue(
			context.Background(),
			arr,
			runtime.NewInt(3),
			runtime.Int(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4,5,3]")
	})
}
