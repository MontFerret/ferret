package arrays_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPop(t *testing.T) {
	Convey("Should return a copy of an array without last element", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		out, err := arrays.Pop(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4]")
	})

	Convey("Should return empty array if a given one is empty", t, func() {
		arr := values.NewArray(0)

		out, err := arrays.Pop(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})
}
