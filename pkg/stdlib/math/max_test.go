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
		So(out, ShouldEqual, runtime.None)

		out, err = math.Max(context.Background(), runtime.NewArray(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}
