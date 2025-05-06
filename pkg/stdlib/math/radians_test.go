package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRadians(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Radians(context.Background(), core.NewInt(180))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3.141592653589793)

		out, err = math.Radians(context.Background(), core.NewFloat(90))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.5707963267948966)

		out, err = math.Radians(context.Background(), core.NewFloat(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)
	})
}
