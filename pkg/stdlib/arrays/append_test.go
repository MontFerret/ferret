package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestAppend(t *testing.T) {
	Convey("Should return a copy of an array", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Append(context.Background(), arr, core.NewInt(6))

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.(collections.Measurable).Length(), ShouldBeGreaterThan, arr.Length())
	})

	Convey("Should ignore non-unique items", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Append(context.Background(), arr, core.NewInt(5), core.True)

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.(collections.Measurable).Length(), ShouldEqual, arr.Length())

		out2, err := arrays.Append(context.Background(), arr, core.NewInt(6), core.True)

		So(err, ShouldBeNil)
		So(out2, ShouldNotEqual, arr)
		So(out2.(collections.Measurable).Length(), ShouldBeGreaterThan, arr.Length())
	})
}
