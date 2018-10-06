package arrays_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReverse(t *testing.T) {
	Convey("Should return a copy of an array with reversed elements", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
			values.NewInt(6),
		)

		out, err := arrays.Reverse(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[6,5,4,3,2,1]")
	})

	Convey("Should return an empty array when there no elements in a source one", t, func() {
		arr := values.NewArray(0)

		out, err := arrays.Reverse(
			context.Background(),
			arr,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})
}
