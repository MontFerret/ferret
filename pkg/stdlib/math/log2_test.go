package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLog2(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Log2(context.Background(), core.NewFloat(1024))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 10)

		out, err = math.Log2(context.Background(), core.NewFloat(8))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Log2(context.Background(), core.NewFloat(0))

		So(err, ShouldBeNil)
		So(core.IsInf(out.(core.Float), -1).Unwrap(), ShouldBeTrue)
	})
}
