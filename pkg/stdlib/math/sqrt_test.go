package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSqrt(t *testing.T) {
	Convey("Should return square value", t, func() {
		out, err := math.Sqrt(context.Background(), values.NewFloat(9))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Sqrt(context.Background(), values.NewInt(2))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.4142135623730951)
	})
}
