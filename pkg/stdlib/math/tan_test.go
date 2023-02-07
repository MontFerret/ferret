package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTan(t *testing.T) {
	Convey("Should return tan value", t, func() {
		out, err := math.Tan(context.Background(), values.NewFloat(10))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, 0.6483608274590867)

		out, err = math.Tan(context.Background(), values.NewInt(5))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, -3.3805150062465854)

		out, err = math.Tan(context.Background(), values.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)
	})
}
