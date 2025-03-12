package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExp2(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Exp2(context.Background(), core.NewFloat(16))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 65536)

		out, err = math.Exp(context.Background(), core.NewFloat(1))

		So(err, ShouldBeNil)
		So(out.Compare(core.NewFloat(2)) == 1, ShouldBeTrue)

		out, err = math.Exp(context.Background(), core.NewFloat(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)
	})
}
