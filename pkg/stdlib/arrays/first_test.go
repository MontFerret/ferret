package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestFirst(t *testing.T) {
	Convey("Should return a first element form a given array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.First(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, 1)
	})

	Convey("Should return NONE if a given array is empty", t, func() {
		arr := values.NewArray(0)

		out, err := arrays.First(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, values.None)
	})
}
