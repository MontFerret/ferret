package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLog2(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Log2(context.Background(), runtime.NewFloat(1024))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 10)

		out, err = math.Log2(context.Background(), runtime.NewFloat(8))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Log2(context.Background(), runtime.NewFloat(0))

		So(err, ShouldBeNil)
		So(runtime.IsInf(out.(runtime.Float), -1).Unwrap(), ShouldBeTrue)
	})
}
