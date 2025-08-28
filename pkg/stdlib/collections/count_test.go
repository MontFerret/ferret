package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

func TestCount(t *testing.T) {
	Convey("When counting an array", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)

		result, err := collections.Count(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(3))
	})

	Convey("When counting an empty array", t, func() {
		arr := runtime.NewArray(0)

		result, err := collections.Count(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(0))
	})

	Convey("When counting a string", t, func() {
		str := runtime.NewString("hello")

		result, err := collections.Count(context.Background(), str)

		So(err, ShouldBeError)
		So(result, ShouldEqual, runtime.ZeroInt)
	})

	Convey("When counting an empty string", t, func() {
		str := runtime.NewString("")

		result, err := collections.Count(context.Background(), str)

		So(err, ShouldBeError)
		So(result, ShouldEqual, runtime.ZeroInt)
	})

	Convey("When counting an object", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewInt(2)),
		)

		result, err := collections.Count(context.Background(), obj)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(2))
	})

	Convey("When counting an empty object", t, func() {
		obj := runtime.NewObject()

		result, err := collections.Count(context.Background(), obj)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(0))
	})

	Convey("When argument is not a collection", t, func() {
		result, err := collections.Count(context.Background(), runtime.NewInt(123))

		So(err, ShouldBeError)
		So(result, ShouldEqual, runtime.ZeroInt)
	})
}
