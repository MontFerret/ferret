package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPush(t *testing.T) {
	Convey("Should create a new array with a new element in the end", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Push(context.Background(), arr, values.NewInt(6))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})

	Convey("Should not add a new element if not unique when uniqueness check is enabled", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Push(
			context.Background(),
			arr,
			values.NewInt(6),
			values.True,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4,5,6]")

		out2, err := arrays.Push(
			context.Background(),
			arr,
			values.NewInt(6),
			values.True,
		)

		So(err, ShouldBeNil)
		So(out2.String(), ShouldEqual, "[1,2,3,4,5,6]")
	})
}
