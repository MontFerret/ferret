package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStandardDeviationPopulation(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.StandardDeviationPopulation(
			context.Background(),
			values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(3),
				values.NewInt(6),
				values.NewInt(5),
				values.NewInt(2),
			),
		)

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, 1.8547236990991407)
	})
}
