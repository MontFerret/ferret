package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestLast(t *testing.T) {
	Convey("Should return a last element form a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Last(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 5)
	})

	Convey("Should return NONE if a given array is empty", t, func() {
		arr := runtime.NewArray(0)

		out, err := arrays.Last(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}
