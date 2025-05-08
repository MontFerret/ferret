package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

func TestPop(t *testing.T) {
	Convey("Should return a copy of an array without last element", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
			runtime.NewInt(4),
			runtime.NewInt(5),
		)

		out, err := arrays.Pop(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[1,2,3,4]")
	})

	Convey("Should return empty array if a given one is empty", t, func() {
		arr := runtime.NewArray(0)

		out, err := arrays.Pop(context.Background(), arr)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "[]")
	})
}
