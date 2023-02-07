package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestNth(t *testing.T) {
	Convey("Should return item by index", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Nth(context.Background(), arr, values.NewInt(1))

		So(err, ShouldBeNil)
		So(out.Compare(values.NewInt(2)), ShouldEqual, 0)
	})

	Convey("Should return None when no value", t, func() {
		arr := values.NewArrayWith()

		out, err := arrays.Nth(context.Background(), arr, values.NewInt(1))

		So(err, ShouldBeNil)
		So(out.Compare(values.None), ShouldEqual, 0)
	})

	Convey("Should return None when passed negative value", t, func() {
		arr := values.NewArrayWith()

		out, err := arrays.Nth(context.Background(), arr, values.NewInt(-1))

		So(err, ShouldBeNil)
		So(out.Compare(values.None), ShouldEqual, 0)
	})
}
