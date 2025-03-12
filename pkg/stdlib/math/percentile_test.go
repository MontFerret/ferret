package math_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPercentile(t *testing.T) {
	Convey("Should return nth percentile value", t, func() {
		out, err := math.Percentile(
			context.Background(),
			internal.NewArrayWith(
				core.NewInt(1),
				core.NewInt(2),
				core.NewInt(3),
				core.NewInt(4),
			),
			core.NewInt(50),
		)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 2)
	})
}
