package math_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPercentileCopyFailure(t *testing.T) {
	for _, method := range []runtime.Value{runtime.String("rank"), runtime.String("interpolation")} {
		t.Run(method.String(), func(t *testing.T) {
			probe := &percentileCopySortable{Value: runtime.True}
			source := &percentileCopyList{List: runtime.NewArrayWith(runtime.Int(3), runtime.Int(1)), copied: probe}
			result, err := math.Percentile(t.Context(), source, runtime.Int(50), method)
			if !errors.Is(err, runtime.ErrInvalidType) {
				t.Fatalf("error = %v", err)
			}

			if value, ok := result.(runtime.Float); !ok || runtime.IsNaN(value) != runtime.True {
				t.Fatalf("result = %v, want NaN", result)
			}

			if source.copies != 1 || probe.sorts != 0 || source.String() != "[3,1]" {
				t.Fatalf("copies = %d, sorts = %d, source = %v", source.copies, probe.sorts, source)
			}
		})
	}
}

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

	Convey("Should support interpolation method", t, func() {
		out, err := math.Percentile(
			context.Background(),
			runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
				runtime.NewInt(4),
			),
			runtime.NewInt(50),
			runtime.NewString("interpolation"),
		)

		So(err, ShouldBeNil)
		So(runtime.Unwrap(out), ShouldEqual, 2.5)
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
