package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPush(t *testing.T) {
	Convey("Should create a new array with a new element in the end", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Push(context.Background(), arr, core.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})

	Convey("Should not add a new element if not unique when uniqueness check is enabled", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Push(
			context.Background(),
			arr,
			core.NewInt(6),
			core.True,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")

		out2, err := arrays.Push(
			context.Background(),
			arr,
			core.NewInt(6),
			core.True,
		)

		So(err, ShouldBeNil)
		So(out2.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})
}
