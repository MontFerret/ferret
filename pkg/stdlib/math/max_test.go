package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMax(t *testing.T) {
	Convey("Should return the largest value", t, func() {
		out, err := math.Max(context.Background(), runtime.NewArrayWith(
			runtime.NewInt(5),
			runtime.NewInt(2),
			runtime.NewInt(9),
			runtime.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 9)

		out, err = math.Max(context.Background(), runtime.NewArrayWith(
			runtime.NewInt(-3),
			runtime.NewInt(-5),
			runtime.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2)

		out, err = math.Max(context.Background(), runtime.NewArrayWith(
			runtime.None,
			runtime.NewInt(-5),
			runtime.False,
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -5)

		out, err = math.Max(context.Background(), runtime.NewArrayWith(
			runtime.None,
			runtime.NewInt(10),
			runtime.False,
			runtime.NewInt(5),
			runtime.NewString("hello"),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 10)

		out, err = math.Max(context.Background(), runtime.NewArray(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}
