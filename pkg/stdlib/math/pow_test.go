package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPow(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Pow(context.Background(), core.NewInt(2), core.NewInt(4))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 16)

		out, err = math.Pow(context.Background(), core.NewInt(5), core.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0.2)

		out, err = math.Pow(context.Background(), core.NewInt(5), core.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)
	})
}
