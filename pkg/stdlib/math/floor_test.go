package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFloor(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Floor(context.Background(), core.NewFloat(2.49))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2)

		out, err = math.Floor(context.Background(), core.NewFloat(2.50))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2)

		out, err = math.Floor(context.Background(), core.NewFloat(-2.50))

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, -3)
	})
}
