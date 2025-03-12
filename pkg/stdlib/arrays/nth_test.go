package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestNth(t *testing.T) {
	Convey("Should return item by index", t, func() {
		arr := internal.NewArrayWith(
			core.NewInt(1),
			core.NewInt(2),
			core.NewInt(3),
			core.NewInt(4),
			core.NewInt(5),
		)

		out, err := arrays.Nth(context.Background(), arr, core.NewInt(1))

		So(err, ShouldBeNil)
		So(out.Compare(core.NewInt(2)), ShouldEqual, 0)
	})

	Convey("Should return None when no value", t, func() {
		arr := internal.NewArrayWith()

		out, err := arrays.Nth(context.Background(), arr, core.NewInt(1))

		So(err, ShouldBeNil)
		So(out.Compare(core.None), ShouldEqual, 0)
	})

	Convey("Should return None when passed negative value", t, func() {
		arr := internal.NewArrayWith()

		out, err := arrays.Nth(context.Background(), arr, core.NewInt(-1))

		So(err, ShouldBeNil)
		So(out.Compare(core.None), ShouldEqual, 0)
	})
}
