package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSqrt(t *testing.T) {
	Convey("Should return square value", t, func() {
		out, err := math.Sqrt(context.Background(), core.NewFloat(9))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Sqrt(context.Background(), core.NewInt(2))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1.4142135623730951)
	})
}
