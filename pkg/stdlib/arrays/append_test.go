package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAppend(t *testing.T) {
	Convey("Should return a copy of an array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := callLegacy("append", context.Background(), arr, runtime.NewInt(6))

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		expected, _ := arr.Length(context.Background())
		actual, _ := out.(runtime.Measurable).Length(context.Background())
		So(actual, ShouldBeGreaterThan, expected)
	})

	Convey("Should ignore non-unique items", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := callLegacy("append", context.Background(), arr, runtime.NewInt(5), runtime.True)

		So(err, ShouldBeNil)
		So(out, ShouldNotPointTo, arr)
		expected, _ := arr.Length(context.Background())
		actual, _ := out.(runtime.Measurable).Length(context.Background())
		So(actual, ShouldEqual, expected)

		out2, err := callLegacy("append", context.Background(), arr, runtime.NewInt(6), runtime.True)

		So(err, ShouldBeNil)
		So(out2, ShouldNotEqual, arr)
		expected, _ = arr.Length(context.Background())
		actual, _ = out2.(runtime.Measurable).Length(context.Background())
		So(actual, ShouldBeGreaterThan, expected)
	})
}
