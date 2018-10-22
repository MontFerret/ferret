package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeRecursive(t *testing.T) {
	Convey("Wrong arguments", t, func() {
		Convey("It should error when 0 arguments", func() {
			actual, err := objects.MergeRecursive(context.Background())

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})

		Convey("It should error when there is not object arguments", func() {
			actual, err := objects.MergeRecursive(context.Background(), values.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)

			actual, err = objects.MergeRecursive(context.Background(),
				values.NewInt(0), values.NewObject(),
			)

			So(err, ShouldBeError)
			So(actual.Compare(values.None), ShouldEqual, 0)
		})
	})

	Convey("Merge single object", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(0)),
		)
		expected := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(0)),
		)

		actual, err := objects.MergeRecursive(context.Background(), obj)

		So(err, ShouldBeNil)
		So(actual.Compare(expected), ShouldEqual, 0)
	})

	Convey("Merge two objects", t, func() {
		Convey("When there are no common keys", t, func() {
			obj1 := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewInt(0)),
			)
			obj2 := values.NewObjectWith(
				values.NewObjectProperty("b", values.NewInt(1)),
			)
			expected := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewInt(0)),
				values.NewObjectProperty("b", values.NewInt(1)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When objects with the same key", func() {
			obj1 := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewInt(0)),
				values.NewObjectProperty("b", values.NewInt(10)),
			)
			obj2 := values.NewObjectWith(
				values.NewObjectProperty("c", values.NewInt(1)),
				values.NewObjectProperty("b", values.NewInt(20)),
			)
			expected := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewInt(1)),
				values.NewObjectProperty("b", values.NewInt(20)),
				values.NewObjectProperty("c", values.NewInt(1)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("Merge two objects with the same keys", func() {
			obj1 := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewInt(0)),
				values.NewObjectProperty("b", values.NewInt(10)),
			)
			obj2 := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewInt(1)),
				values.NewObjectProperty("b", values.NewInt(20)),
			)
			expected := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewInt(1)),
				values.NewObjectProperty("b", values.NewInt(20)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When there are nested arrays", t, func() {

		})

		Convey("When there are nested objects", t, func() {

		})

		Convey("When there are nested objects and arrays", t, func() {

		})

		Convey("When there are nested simple and complex types", t, func() {

		})
	})

	Convey("Merged object should be independent of source objects", t, func() {
		Convey("When simple types", func() {

		})

		Convey("When arrays", func() {

		})

		Convey("When objects", func() {

		})

		Convey("When arrays and objects", func() {

		})

		Convey("When complex and simple types", func() {

		})
	})
}
