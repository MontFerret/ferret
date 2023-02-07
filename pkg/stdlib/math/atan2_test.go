package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAtan2(t *testing.T) {
	Convey("Should return tangent value", t, func() {
		out, err := math.Atan2(context.Background(), values.NewInt(0), values.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		out, err = math.Atan2(context.Background(), values.NewInt(1), values.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.5707963267948966)

		out, err = math.Atan2(context.Background(), values.NewInt(1), values.NewInt(1))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, 0.7853981633974483)

		out, err = math.Atan2(context.Background(), values.NewInt(-10), values.NewInt(20))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, -0.4636476090008061)
	})
}
