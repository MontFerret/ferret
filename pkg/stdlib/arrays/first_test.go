package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestFirst(t *testing.T) {
	Convey("Should return a first element form a given array", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.First(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)
	})

	Convey("Should return NONE if a given array is empty", t, func() {
		arr := internal.NewArray(0)

		out, err := arrays.First(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, core.None)
	})
}
