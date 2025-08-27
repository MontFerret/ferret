package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSqrt(t *testing.T) {
	Convey("Should return square value", t, func() {
		out, err := math.Sqrt(context.Background(), runtime.NewFloat(9))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Sqrt(context.Background(), runtime.NewInt(2))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.4142135623730951)

		out, err = math.Sqrt(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		// Test negative numbers
		out, err = math.Sqrt(context.Background(), runtime.NewFloat(-4))

		So(err, ShouldBeNil)
		So(runtime.IsNaN(out.(runtime.Float)).Unwrap(), ShouldBeTrue)
	})

	Convey("Should return error when value is not numeric", t, func() {
		out, err := math.Sqrt(context.Background(), runtime.NewString("invalid"))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)

		out, err = math.Sqrt(context.Background(), runtime.None)

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}
