package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExp(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.Exp(context.Background(), core.NewFloat(1))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2.718281828459045)

		out, err = math.Exp(context.Background(), core.NewFloat(10))

		So(err, ShouldBeNil)
		So(out.Compare(core.NewFloat(22026.46579480671)) == 1, ShouldBeTrue)

		out, err = math.Exp(context.Background(), core.NewFloat(0))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)
	})
}
