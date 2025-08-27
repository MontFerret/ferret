package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAtan2(t *testing.T) {
	Convey("Should return tangent value", t, func() {
		out, err := math.Atan2(context.Background(), runtime.NewInt(0), runtime.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		out, err = math.Atan2(context.Background(), runtime.NewInt(1), runtime.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.5707963267948966)

		out, err = math.Atan2(context.Background(), runtime.NewInt(1), runtime.NewInt(1))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, 0.7853981633974483)

		out, err = math.Atan2(context.Background(), runtime.NewInt(-10), runtime.NewInt(20))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, -0.4636476090008061)
	})
}
