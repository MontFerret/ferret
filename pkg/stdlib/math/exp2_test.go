package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExp2(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Exp2(context.Background(), values.NewFloat(16))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 65536)

		out, err = math.Exp(context.Background(), values.NewFloat(1))

		So(err, ShouldBeNil)
		So(out.Compare(values.NewFloat(2)) == 1, ShouldBeTrue)

		out, err = math.Exp(context.Background(), values.NewFloat(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)
	})
}
