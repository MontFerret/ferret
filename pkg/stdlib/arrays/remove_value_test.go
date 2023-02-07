package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestRemoveValue(t *testing.T) {
	Convey("Should return a copy of an array without given element(s)", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(3),
		)

		out, err := arrays.RemoveValue(context.Background(), arr, values.NewInt(3))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4]")
	})

	Convey("Should return a copy of an array without given element(s) with limit", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(3),
			values.NewInt(5),
			values.NewInt(3),
		)

		out, err := arrays.RemoveValue(
			context.Background(),
			arr,
			values.NewInt(3),
			values.Int(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,4,5,3]")
	})
}
