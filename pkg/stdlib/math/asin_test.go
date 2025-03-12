package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAsin(t *testing.T) {
	Convey("Should return arcsine value", t, func() {
		out, err := math.Asin(context.Background(), core.NewInt(1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.5707963267948966)

		out, err = math.Asin(context.Background(), core.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		out, err = math.Asin(context.Background(), core.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -1.5707963267948966)

		out, err = math.Asin(context.Background(), core.NewInt(2))

		So(err, ShouldBeNil)
		So(core.IsNaN(out.(core.Float)), ShouldEqual, true)
	})
}
