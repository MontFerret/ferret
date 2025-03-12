package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMedian(t *testing.T) {
	Convey("Should return median value", t, func() {
		out, err := math.Median(context.Background(), internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2)

		out, err = math.Average(context.Background(), internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2.5)

		out, err = math.Average(context.Background(), internal.NewArrayWith(
			core.NewInt(2),
			core.NewInt(1),
			core.NewInt(4),
			core.NewInt(3),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2.5)

		out, err = math.Average(context.Background(), internal.NewArrayWith(
			core.None,
			core.NewInt(-5),
			core.False,
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, core.None)
	})
}
