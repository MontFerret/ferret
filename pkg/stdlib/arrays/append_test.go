package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestAppend(t *testing.T) {
	Convey("Should return a copy of an array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Append(context.Background(), arr, values.NewInt(6))

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.(collections.Measurable).Length(), ShouldBeGreaterThan, arr.Length())
	})

	Convey("Should ignore non-unique items", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Append(context.Background(), arr, values.NewInt(5), values.True)

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.(collections.Measurable).Length(), ShouldEqual, arr.Length())

		out2, err := arrays.Append(context.Background(), arr, values.NewInt(6), values.True)

		So(err, ShouldBeNil)
		So(out2, ShouldNotEqual, arr)
		So(out2.(collections.Measurable).Length(), ShouldBeGreaterThan, arr.Length())
	})
}
