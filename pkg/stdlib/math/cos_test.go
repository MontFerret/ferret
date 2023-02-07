package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCos(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Cos(context.Background(), values.NewFloat(1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0.5403023058681398)

		out, err = math.Cos(context.Background(), values.NewFloat(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)

		out, err = math.Cos(context.Background(), values.NewFloat(-3.141592653589783))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, -1)
	})
}
