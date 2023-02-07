package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSin(t *testing.T) {
	Convey("Should return sin value", t, func() {
		out, err := math.Sin(context.Background(), values.NewFloat(3.141592653589783/2))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)

		out, err = math.Sin(context.Background(), values.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		out, err = math.Sin(context.Background(), values.NewFloat(-3.141592653589783/2))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -1)

		v, _ := math.Radians(context.Background(), values.NewInt(270))

		out, err = math.Sin(context.Background(), v)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -1)
	})
}
