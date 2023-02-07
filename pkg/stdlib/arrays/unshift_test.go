package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestUnshift(t *testing.T) {
	Convey("Should return a copy of an array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Unshift(context.Background(), arr, values.NewInt(0))

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")
	})

	Convey("Should ignore non-unique items", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Unshift(
			context.Background(),
			arr,
			values.NewInt(0),
			values.True,
		)

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")

		out2, err := arrays.Unshift(
			context.Background(),
			arr,
			values.NewInt(0),
			values.True,
		)

		So(err, ShouldBeNil)
		So(out2, ShouldNotEqual, arr)
		So(out.String(), ShouldEqual, "[0,1,2,3,4,5]")
	})
}
