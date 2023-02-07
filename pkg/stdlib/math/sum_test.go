package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSum(t *testing.T) {
	Convey("Should return sum of values", t, func() {
		out, err := math.Sum(context.Background(), values.NewArrayWith(
			values.NewInt(5),
			values.NewInt(2),
			values.NewInt(9),
			values.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 18)

		out, err = math.Sum(context.Background(), values.NewArrayWith(
			values.NewInt(-3),
			values.NewInt(-5),
			values.NewInt(2),
		))

		So(err, ShouldBeNil)
		So(out, ShouldEqual, -6)
	})
}
