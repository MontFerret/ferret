package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMedian(t *testing.T) {
	Convey("Should return median value", t, func() {
		out, err := math.Median(context.Background(), runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2)

		out, err = math.Median(context.Background(), runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2.5)

		out, err = math.Median(context.Background(), runtime.NewArrayWith(
			runtime.NewInt(2),
			runtime.NewInt(1),
			runtime.NewInt(4),
			runtime.NewInt(3),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2.5)

		out, err = math.Median(context.Background(), runtime.NewArrayWith(
			runtime.None,
			runtime.NewInt(-5),
			runtime.False,
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}
