package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValues(t *testing.T) {
	Convey("Invalid arguments", t, func() {
		Convey("When there is no arguments", func() {
			actual, err := objects.Values(context.Background())

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("When 2 arguments", func() {
			obj := values.NewObjectWith(
				values.NewObjectProperty("k1", values.NewInt(0)),
				values.NewObjectProperty("k2", values.NewInt(1)),
			)

			actual, err := objects.Values(context.Background(), obj, obj)

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)

			actual, err = objects.Values(context.Background(), obj, values.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("When there is not object argument", func() {
			actual, err := objects.Values(context.Background(), values.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})
	})

	Convey("When simple type attributes (same type)", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("k1", values.NewInt(0)),
			values.NewObjectProperty("k2", values.NewInt(1)),
		)
		expected := values.NewArrayWith(
			values.NewInt(0), values.NewInt(1),
		).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*values.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When simple type attributes (different types)", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("k1", values.NewInt(0)),
			values.NewObjectProperty("k2", values.NewString("v2")),
		)
		expected := values.NewArrayWith(
			values.NewInt(0), values.NewString("v2"),
		).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*values.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (array)", t, func() {
		arr1 := values.NewArrayWith(
			values.NewInt(0), values.NewInt(1),
		)
		arr2 := values.NewArrayWith(
			values.NewInt(2), values.NewInt(3),
		)
		obj := values.NewObjectWith(
			values.NewObjectProperty("k1", arr1),
			values.NewObjectProperty("k2", arr2),
		)
		expected := values.NewArrayWith(arr1, arr2).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*values.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (object)", t, func() {
		obj1 := values.NewObjectWith(
			values.NewObjectProperty("int0", values.NewInt(0)),
		)
		obj2 := values.NewObjectWith(
			values.NewObjectProperty("int1", values.NewInt(1)),
		)
		obj := values.NewObjectWith(
			values.NewObjectProperty("k1", obj1),
			values.NewObjectProperty("k2", obj2),
		)
		expected := values.NewArrayWith(obj1, obj2).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*values.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (object and array)", t, func() {
		obj1 := values.NewObjectWith(
			values.NewObjectProperty("k1", values.NewInt(0)),
		)
		arr1 := values.NewArrayWith(
			values.NewInt(0), values.NewInt(1),
		)
		obj := values.NewObjectWith(
			values.NewObjectProperty("obj", obj1),
			values.NewObjectProperty("arr", arr1),
		)
		expected := values.NewArrayWith(obj1, arr1).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*values.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When both type attributes", t, func() {
		obj1 := values.NewObjectWith(
			values.NewObjectProperty("k1", values.NewInt(0)),
		)
		arr1 := values.NewArrayWith(
			values.NewInt(0), values.NewInt(1),
		)
		int1 := values.NewInt(0)
		obj := values.NewObjectWith(
			values.NewObjectProperty("obj", obj1),
			values.NewObjectProperty("arr", arr1),
			values.NewObjectProperty("int", int1),
		)
		expected := values.NewArrayWith(obj1, arr1, int1).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*values.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("Result is independent on the source object (array)", t, func() {
		arr := values.NewArrayWith(values.NewInt(0))
		obj := values.NewObjectWith(
			values.NewObjectProperty("arr", arr),
		)
		expected := values.NewArrayWith(values.NewInt(0))

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*values.Array).Sort()

		So(err, ShouldBeNil)

		arr.Push(values.NewInt(1))

		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("Result is independent on the source object (object)", t, func() {
		nested := values.NewObjectWith(
			values.NewObjectProperty("int", values.NewInt(0)),
		)
		obj := values.NewObjectWith(
			values.NewObjectProperty("nested", nested),
		)
		expected := values.NewArrayWith(
			values.NewObjectWith(
				values.NewObjectProperty("int", values.NewInt(0)),
			),
		)

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*values.Array).Sort()

		So(err, ShouldBeNil)

		nested.Set("new", values.NewInt(1))

		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})
}
