package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestLast_Basic(t *testing.T) {
	ctx := context.Background()

	Convey("Should return a last element form a given array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Last(ctx, arr)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 5)
	})

	Convey("Should return NONE if a given array is empty", t, func() {
		arr := runtime.NewArray(0)

		out, err := arrays.Last(ctx, arr)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}

func TestLast_ArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Should properly validate arguments", t, func() {
		nonArray := runtime.NewString("not an array")

		Convey("Should return error for non-array input", func() {
			_, err := arrays.Last(ctx, nonArray)
			So(err, ShouldNotBeNil)
		})

		Convey("Should work correctly with valid array", func() {
			arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
			out, err := arrays.Last(ctx, arr)
			So(err, ShouldBeNil)
			So(out.(runtime.Comparable).Compare(runtime.NewInt(2)), ShouldEqual, 0)
		})
	})
}

func TestLast_SpecialValues(t *testing.T) {
	ctx := context.Background()

	Convey("Should handle arrays with only None values", t, func() {
		noneOnlyArr := runtime.NewArrayWith(runtime.None, runtime.None, runtime.None)
		out, err := arrays.Last(ctx, noneOnlyArr)
		So(err, ShouldBeNil)
		So(out, ShouldEqual, runtime.None)
	})
}
