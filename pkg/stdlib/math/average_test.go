package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAverage(t *testing.T) {
	Convey("Should return average value", t, func() {
		out, err := math.Average(context.Background(), internal.NewArrayWith(
			core.NewInt(5),
			core.NewInt(2),
			core.NewInt(9),
			core.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 4.5)

		out, err = math.Average(context.Background(), internal.NewArrayWith(
			core.NewInt(-3),
			core.NewInt(-5),
			core.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -2)

		out, err = math.Average(context.Background(), internal.NewArrayWith(
			core.None,
			core.NewInt(-5),
			core.False,
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, core.None)

		out, err = math.Average(context.Background(), internal.NewArray(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, core.None)
	})
}
