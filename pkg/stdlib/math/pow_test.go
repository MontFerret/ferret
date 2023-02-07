package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPow(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Pow(context.Background(), values.NewInt(2), values.NewInt(4))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 16)

		out, err = math.Pow(context.Background(), values.NewInt(5), values.NewInt(-1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 0.2)

		out, err = math.Pow(context.Background(), values.NewInt(5), values.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)
	})
}
