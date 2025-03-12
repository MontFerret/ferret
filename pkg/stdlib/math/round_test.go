package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRound(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Round(context.Background(), core.NewFloat(2.49))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2)

		out, err = math.Round(context.Background(), core.NewFloat(2.50))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 3)

		out, err = math.Ceil(context.Background(), core.NewFloat(-2.50))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, -2)

		out, err = math.Ceil(context.Background(), core.NewFloat(-2.49))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, -2)
	})
}
