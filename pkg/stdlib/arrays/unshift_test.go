package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnshift(t *testing.T) {
	Convey("Should return a copy of an array", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Unshift(context.Background(), arr, core.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")
	})

	Convey("Should ignore non-unique items", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Unshift(
			context.Background(),
			arr,
			core.NewInt(0),
			core.True,
		)

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")

		out2, err := arrays.Unshift(
			context.Background(),
			arr,
			core.NewInt(0),
			core.True,
		)

		So(err, ShouldBeNil)
		So(out2, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")
	})
}
