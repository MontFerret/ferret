package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLog10(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Log10(context.Background(), runtime.NewFloat(10000))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 4)

		out, err = math.Log10(context.Background(), runtime.NewFloat(10))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)

		out, err = math.Log10(context.Background(), runtime.NewFloat(0))

		So(err, ShouldBeNil)
		So(runtime.IsInf(out.(runtime.Float), -1).Unwrap(), ShouldBeTrue)
	})
}
