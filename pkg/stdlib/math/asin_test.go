package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAsin(t *testing.T) {
	Convey("Should return arcsine value", t, func() {
		out, err := math.Asin(context.Background(), runtime.NewInt(1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.5707963267948966)

		out, err = math.Asin(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		out, err = math.Asin(context.Background(), runtime.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -1.5707963267948966)

		out, err = math.Asin(context.Background(), runtime.NewInt(2))

		So(err, ShouldBeNil)
		So(runtime.IsNaN(out.(runtime.Float)).Unwrap(), ShouldBeTrue)
	})
}
