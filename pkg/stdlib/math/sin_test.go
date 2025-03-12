package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSin(t *testing.T) {
	Convey("Should return sin value", t, func() {
		out, err := math.Sin(context.Background(), core.NewFloat(3.141592653589783/2))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)

		out, err = math.Sin(context.Background(), core.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		out, err = math.Sin(context.Background(), core.NewFloat(-3.141592653589783/2))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -1)

		v, _ := math.Radians(context.Background(), core.NewInt(270))

		out, err = math.Sin(context.Background(), v)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -1)
	})
}
