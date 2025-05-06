package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveValue(t *testing.T) {
	Convey("Should return a copy of an array without given element(s)", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(3),
		)

		out, err := arrays.RemoveValue(context.Background(), arr, core.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4]")
	})

	Convey("Should return a copy of an array without given element(s) with limit", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(3),
			core.NewInt(5),
			core.NewInt(3),
		)

		out, err := arrays.RemoveValue(
			context.Background(),
			arr,
			core.NewInt(3),
			core.Int(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4,5,3]")
	})
}
