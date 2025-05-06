package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAtan(t *testing.T) {
	Convey("Should return arctangent value", t, func() {
		out, err := math.Atan(context.Background(), core.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -0.7853981633974483)

		out, err = math.Atan(context.Background(), core.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		out, err = math.Atan(context.Background(), core.NewInt(10))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, 1.4711276743037345)
	})
}
