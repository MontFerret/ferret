package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAverage(t *testing.T) {
	Convey("Should return average value", t, func() {
		out, err := math.Average(context.Background(), runtime.NewArrayWith(
			runtime.NewInt(5),
			runtime.NewInt(2),
			runtime.NewInt(9),
			runtime.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 4.5)

		out, err = math.Average(context.Background(), runtime.NewArrayWith(
			runtime.NewInt(-3),
			runtime.NewInt(-5),
			runtime.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -2)

		out, err = math.Average(context.Background(), runtime.NewArrayWith(
			runtime.None,
			runtime.NewInt(-5),
			runtime.False,
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.Float(-5))

		out, err = math.Average(context.Background(), runtime.NewArray(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.Float(0))
	})

	Convey("Should ignore nulls and compute correct average", t, func() {
		out, err := math.Average(context.Background(), runtime.NewArrayWith(
			runtime.None,
			runtime.NewInt(20),
			runtime.NewInt(0),
			runtime.None,
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.Float(10.0))
	})
}
