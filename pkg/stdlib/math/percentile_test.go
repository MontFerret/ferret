package math_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPercentile(t *testing.T) {
	Convey("Should return nth percentile value", t, func() {
		out, err := math.Percentile(
			context.Background(),
			runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
				runtime.NewInt(4),
			),
			runtime.NewInt(50),
		)

		So(err, ShouldBeNil)
		So(runtime.Unwrap(out), ShouldEqual, 2)

		// Test with different percentile
		out, err = math.Percentile(
			context.Background(),
			runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
				runtime.NewInt(4),
				runtime.NewInt(5),
			),
			runtime.NewInt(80),
		)

		So(err, ShouldBeNil)
		So(runtime.Unwrap(out), ShouldEqual, 4)
	})

	Convey("Should return error for invalid arguments", t, func() {
		// Non-array first argument
		out, err := math.Percentile(context.Background(), runtime.NewInt(1), runtime.NewInt(50))

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)

		// Non-numeric percentile
		out, err = math.Percentile(
			context.Background(),
			runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
			runtime.NewString("invalid"),
		)

		So(err, ShouldNotBeNil)
		So(out, ShouldEqual, runtime.None)

		// Empty array returns NaN, not error
		out, err = math.Percentile(
			context.Background(),
			runtime.NewArray(0),
			runtime.NewInt(50),
		)

		So(err, ShouldBeNil)
		So(runtime.Unwrap(runtime.IsNaN(out.(runtime.Float))), ShouldBeTrue)

		// Invalid percentile (outside range)
		_, err = math.Percentile(
			context.Background(),
			runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
			runtime.NewInt(0),
		)

		So(err, ShouldNotBeNil)
	})
}
