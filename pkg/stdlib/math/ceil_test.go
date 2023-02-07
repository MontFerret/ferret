package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCeil(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Ceil(context.Background(), values.NewFloat(2.49))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Ceil(context.Background(), values.NewFloat(2.50))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Ceil(context.Background(), values.NewFloat(-2.50))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, -2)
	})
}
