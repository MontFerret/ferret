package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValues(t *testing.T) {
	Convey("Invalid arguments", t, func() {
		Convey("When there is no arguments", func() {
			actual, err := objects.Values(context.Background())

			So(err, ShouldBeError)
			So(actual.Compare(core.None), ShouldEqual, 0)
		})

		Convey("When 2 arguments", func() {
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("k1", core.NewInt(0)),
				internal.NewObjectProperty("k2", core.NewInt(1)),
			)

			actual, err := objects.Values(context.Background(), obj, obj)

			So(err, ShouldBeError)
			So(actual.Compare(core.None), ShouldEqual, 0)

			actual, err = objects.Values(context.Background(), obj, core.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(core.None), ShouldEqual, 0)
		})

		Convey("When there is not object argument", func() {
			actual, err := objects.Values(context.Background(), core.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(core.None), ShouldEqual, 0)
		})
	})

	Convey("When simple type attributes (same type)", t, func() {
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("k1", core.NewInt(0)),
			internal.NewObjectProperty("k2", core.NewInt(1)),
		)
		expected := internal.NewArrayWith(
			core.NewInt(0), core.NewInt(1),
		).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*internal.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When simple type attributes (different types)", t, func() {
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("k1", core.NewInt(0)),
			internal.NewObjectProperty("k2", core.NewString("v2")),
		)
		expected := internal.NewArrayWith(
			core.NewInt(0), core.NewString("v2"),
		).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*internal.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (array)", t, func() {
		arr1 := internal.NewArrayWith(
			core.NewInt(0), core.NewInt(1),
		)
		arr2 := internal.NewArrayWith(
			core.NewInt(2), core.NewInt(3),
		)
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("k1", arr1),
			internal.NewObjectProperty("k2", arr2),
		)
		expected := internal.NewArrayWith(arr1, arr2).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*internal.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (object)", t, func() {
		obj1 := internal.NewObjectWith(
			internal.NewObjectProperty("int0", core.NewInt(0)),
		)
		obj2 := internal.NewObjectWith(
			internal.NewObjectProperty("int1", core.NewInt(1)),
		)
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("k1", obj1),
			internal.NewObjectProperty("k2", obj2),
		)
		expected := internal.NewArrayWith(obj1, obj2).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*internal.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (object and array)", t, func() {
		obj1 := internal.NewObjectWith(
			internal.NewObjectProperty("k1", core.NewInt(0)),
		)
		arr1 := internal.NewArrayWith(
			core.NewInt(0), core.NewInt(1),
		)
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("obj", obj1),
			internal.NewObjectProperty("arr", arr1),
		)
		expected := internal.NewArrayWith(obj1, arr1).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*internal.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When both type attributes", t, func() {
		obj1 := internal.NewObjectWith(
			internal.NewObjectProperty("k1", core.NewInt(0)),
		)
		arr1 := internal.NewArrayWith(
			core.NewInt(0), core.NewInt(1),
		)
		int1 := core.NewInt(0)
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("obj", obj1),
			internal.NewObjectProperty("arr", arr1),
			internal.NewObjectProperty("int", int1),
		)
		expected := internal.NewArrayWith(obj1, arr1, int1).Sort()

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*internal.Array).Sort()

		So(err, ShouldBeNil)
		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("Result is independent on the source object (array)", t, func() {
		arr := internal.NewArrayWith(core.NewInt(0))
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("arr", arr),
		)
		expected := internal.NewArrayWith(
			internal.NewArrayWith(
				core.NewInt(0),
			),
		)

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*internal.Array).Sort()

		So(err, ShouldBeNil)

		arr.Push(core.NewInt(1))

		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("Result is independent on the source object (object)", t, func() {
		nested := internal.NewObjectWith(
			internal.NewObjectProperty("int", core.NewInt(0)),
		)
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("nested", nested),
		)
		expected := internal.NewArrayWith(
			internal.NewObjectWith(
				internal.NewObjectProperty("int", core.NewInt(0)),
			),
		)

		actual, err := objects.Values(context.Background(), obj)
		actualSorted := actual.(*internal.Array).Sort()

		So(err, ShouldBeNil)

		nested.Set("new", core.NewInt(1))

		So(actualSorted.Compare(expected), ShouldEqual, 0)
	})
}

func TestValuesStress(t *testing.T) {
	Convey("Stress", t, func() {
		for i := 0; i < 100; i++ {
			obj1 := internal.NewObjectWith(
				internal.NewObjectProperty("int0", core.NewInt(0)),
			)
			obj2 := internal.NewObjectWith(
				internal.NewObjectProperty("int1", core.NewInt(1)),
			)
			obj := internal.NewObjectWith(
				internal.NewObjectProperty("k1", obj1),
				internal.NewObjectProperty("k2", obj2),
			)
			expected := internal.NewArrayWith(obj2, obj1).Sort()

			actual, err := objects.Values(context.Background(), obj)
			actualSorted := actual.(*internal.Array).Sort()

			So(err, ShouldBeNil)
			So(actualSorted.Length(), ShouldEqual, expected.Length())
			So(actualSorted.Compare(expected), ShouldEqual, 0)
		}
	})
}
