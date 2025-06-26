package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStandardDeviationSample(t *testing.T) {
	Convey("Should return a value", t, func() {
		out, err := math.StandardDeviationSample(
			context.Background(),
			runtime.NewArrayWith(
				core.NewInt(1),
				core.NewInt(3),
				core.NewInt(6),
				core.NewInt(5),
				core.NewInt(2),
			),
		)

		So(err, ShouldBeNil)
		So(out.Unwrap(), ShouldEqual, 2.073644135332772)
	})
}
