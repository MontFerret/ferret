package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLog(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Log(context.Background(), core.NewFloat(2.718281828459045))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)

		out, err = math.Log(context.Background(), core.NewFloat(10))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2.302585092994046)

		out, err = math.Log(context.Background(), core.NewFloat(0))

		So(err, ShouldBeNil)
		So(core.IsInf(out.(core.Float), -1).Unwrap(), ShouldBeTrue)
	})
}
