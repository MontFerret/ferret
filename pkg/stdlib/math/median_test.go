package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMedian(t *testing.T) {
	Convey("Should return median value", t, func() {
		out, err := math.Median(context.Background(), values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2)

		out, err = math.Average(context.Background(), values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2.5)

		out, err = math.Average(context.Background(), values.NewArrayWith(
			values.NewInt(2),
			values.NewInt(1),
			values.NewInt(4),
			values.NewInt(3),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2.5)

		out, err = math.Average(context.Background(), values.NewArrayWith(
			values.None,
			values.NewInt(-5),
			values.False,
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, values.None)
	})
}
