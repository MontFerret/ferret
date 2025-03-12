package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAbs(t *testing.T) {
	Convey("Should return absolute value", t, func() {
		Convey("When value is int", func() {
			out, err := math.Abs(context.Background(), core.NewInt(-5))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, 5)

			out, err = math.Abs(context.Background(), core.NewInt(3))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, 3)
		})

		Convey("When value is float", func() {
			out, err := math.Abs(context.Background(), core.NewFloat(-5))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, 5)

			out, err = math.Abs(context.Background(), core.NewFloat(5.1))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, 5.1)
		})
	})
}
