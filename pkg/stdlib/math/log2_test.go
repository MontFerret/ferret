package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLog2(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Log2(context.Background(), values.NewFloat(1024))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 10)

		out, err = math.Log2(context.Background(), values.NewFloat(8))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Log2(context.Background(), values.NewFloat(0))

		So(err, ShouldBeNil)
		So(values.IsInf(out.(values.Float), -1).Unwrap(), ShouldBeTrue)
	})
}
