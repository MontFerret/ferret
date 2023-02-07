package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAsin(t *testing.T) {
	Convey("Should return arcsine value", t, func() {
		out, err := math.Asin(context.Background(), values.NewInt(1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.5707963267948966)

		out, err = math.Asin(context.Background(), values.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0)

		out, err = math.Asin(context.Background(), values.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -1.5707963267948966)

		out, err = math.Asin(context.Background(), values.NewInt(2))

		So(err, ShouldBeNil)
		So(values.IsNaN(out.(values.Float)), ShouldEqual, true)
	})
}
