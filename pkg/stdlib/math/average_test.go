package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAverage(t *testing.T) {
	Convey("Should return average value", t, func() {
		out, err := math.Average(context.Background(), values.NewArrayWith(
			values.NewInt(5),
			values.NewInt(2),
			values.NewInt(9),
			values.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 4.5)

		out, err = math.Average(context.Background(), values.NewArrayWith(
			values.NewInt(-3),
			values.NewInt(-5),
			values.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -2)

		out, err = math.Average(context.Background(), values.NewArrayWith(
			values.None,
			values.NewInt(-5),
			values.False,
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, values.None)

		out, err = math.Average(context.Background(), values.NewArray(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, values.None)
	})
}
