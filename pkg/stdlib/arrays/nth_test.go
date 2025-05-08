package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestNth(t *testing.T) {
	Convey("Should return item by index", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Nth(context.Background(), arr, runtime.NewInt(1))

		So(err, ShouldBeNil)
		So(out.(runtime.Comparable).Compare(runtime.NewInt(2)), ShouldEqual, 0)
	})

	Convey("Should return None when no value", t, func() {
		arr := runtime.NewArrayWith()

		out, err := arrays.Nth(context.Background(), arr, runtime.NewInt(1))

		So(err, ShouldBeNil)
		So(out, ShouldPointTo, runtime.None)
	})

	Convey("Should return None when passed negative value", t, func() {
		arr := runtime.NewArrayWith()

		out, err := arrays.Nth(context.Background(), arr, runtime.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldPointTo, runtime.None)
	})
}
