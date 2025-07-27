package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSum(t *testing.T) {
	Convey("Should return sum of values", t, func() {
		out, err := math.Sum(context.Background(), runtime.NewArrayWith(
			core.NewInt(5),
			core.NewInt(2),
			core.NewInt(9),
			core.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 18)

		out, err = math.Sum(context.Background(), runtime.NewArrayWith(
			core.NewInt(-3),
			core.NewInt(-5),
			core.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -6)
	})

	Convey("Should ignore none number-values", t, func() {
		out, err := math.Sum(context.Background(), runtime.NewArrayWith(
			core.None,
			core.NewInt(2),
			core.True,
			core.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 4)
	})
}
