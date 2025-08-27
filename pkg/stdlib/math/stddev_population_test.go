package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStandardDeviationPopulation(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.StandardDeviationPopulation(
			context.Background(),
			runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(3),
				runtime.NewInt(6),
				runtime.NewInt(5),
				runtime.NewInt(2),
			),
		)

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, 1.8547236990991407)
	})
}
