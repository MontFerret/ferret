package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

func TestCountDistinct(t *testing.T) {
	Convey("When counting distinct elements in array with duplicates", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(1),
			runtime.NewInt(3),
			runtime.NewInt(2),
		)

		result, err := collections.CountDistinct(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(3))
	})

	Convey("When counting distinct elements in array without duplicates", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewInt(2),
			runtime.NewInt(3),
		)

		result, err := collections.CountDistinct(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(3))
	})

	Convey("When counting distinct elements in empty array", t, func() {
		arr := runtime.NewArray(0)

		result, err := collections.CountDistinct(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(0))
	})

	Convey("When counting distinct elements with different types", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(1),
			runtime.NewString("1"),
			runtime.NewInt(1),
			runtime.NewString("hello"),
			runtime.NewString("1"),
		)

		result, err := collections.CountDistinct(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(3)) // 1 (int), "1" (string), "hello" (string)
	})

	Convey("When counting distinct elements in object", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewInt(2)),
			runtime.NewObjectProperty("c", runtime.NewInt(1)), // duplicate value
		)

		result, err := collections.CountDistinct(context.Background(), obj)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(2)) // values 1 and 2
	})

	Convey("When counting distinct elements in empty object", t, func() {
		obj := runtime.NewObject()

		result, err := collections.CountDistinct(context.Background(), obj)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(0))
	})

	Convey("When argument is not a collection", t, func() {
		result, err := collections.CountDistinct(context.Background(), runtime.NewInt(123))

		So(err, ShouldBeError)
		So(result, ShouldEqual, runtime.ZeroInt)
	})

	Convey("When counting distinct elements with all same values", t, func() {
		arr := runtime.NewArrayWith(
			runtime.NewInt(5),
			runtime.NewInt(5),
			runtime.NewInt(5),
			runtime.NewInt(5),
		)

		result, err := collections.CountDistinct(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(1))
	})

	Convey("When counting distinct with complex objects", t, func() {
		obj1 := runtime.NewObjectWith(
			runtime.NewObjectProperty("x", runtime.NewInt(1)),
		)
		obj2 := runtime.NewObjectWith(
			runtime.NewObjectProperty("x", runtime.NewInt(1)),
		)
		obj3 := runtime.NewObjectWith(
			runtime.NewObjectProperty("x", runtime.NewInt(2)),
		)

		arr := runtime.NewArrayWith(obj1, obj2, obj3)

		result, err := collections.CountDistinct(context.Background(), arr)

		So(err, ShouldBeNil)
		// Objects with same content should be considered the same
		So(result, ShouldEqual, runtime.NewInt(2))
	})

	Convey("When counting distinct with None values", t, func() {
		arr := runtime.NewArrayWith(
			runtime.None,
			runtime.NewInt(1),
			runtime.None,
			runtime.NewString("test"),
		)

		result, err := collections.CountDistinct(context.Background(), arr)

		So(err, ShouldBeNil)
		So(result, ShouldEqual, runtime.NewInt(3)) // None, 1, "test"
	})
}