package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDegrees(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Degrees(context.Background(), core.NewFloat(0.7853981633974483))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 45)

		out, err = math.Degrees(context.Background(), core.NewFloat(3.141592653589793))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 180)
	})
}
