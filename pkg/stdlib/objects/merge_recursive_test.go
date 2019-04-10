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
		Convey("When there are no common keys", func() {
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
				values.NewObjectProperty("a", values.NewInt(0)),
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

		Convey("When there are nested arrays", func() {
			obj1 := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewArrayWith(
					values.NewInt(1), values.NewInt(2),
				)),
			)
			obj2 := values.NewObjectWith(
				values.NewObjectProperty("b", values.NewArrayWith(
					values.NewInt(1), values.NewInt(2),
				)),
			)
			expected := values.NewObjectWith(
				values.NewObjectProperty("a", values.NewArrayWith(
					values.NewInt(1), values.NewInt(2),
				)),
				values.NewObjectProperty("b", values.NewArrayWith(
					values.NewInt(1), values.NewInt(2),
				)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When there are nested objects (example from ArangoDB doc)", func() {
			// { "user-1": { "name": "Jane", "livesIn": { "city": "LA" } } }
			obj1 := values.NewObjectWith(
				values.NewObjectProperty(
					"user-1", values.NewObjectWith(
						values.NewObjectProperty(
							"name", values.NewString("Jane"),
						),
						values.NewObjectProperty(
							"livesIn", values.NewObjectWith(
								values.NewObjectProperty(
									"city", values.NewString("LA"),
								),
							),
						),
					),
				),
			)
			// { "user-1": { "age": 42, "livesIn": { "state": "CA" } } }
			obj2 := values.NewObjectWith(
				values.NewObjectProperty(
					"user-1", values.NewObjectWith(
						values.NewObjectProperty(
							"age", values.NewInt(42),
						),
						values.NewObjectProperty(
							"livesIn", values.NewObjectWith(
								values.NewObjectProperty(
									"state", values.NewString("CA"),
								),
							),
						),
					),
				),
			)
			// { "user-1": { "age": 42, "livesIn": { "city": "LA", "state": "CA" }, "name": "Jane" } }
			expected := values.NewObjectWith(
				values.NewObjectProperty(
					"user-1", values.NewObjectWith(
						values.NewObjectProperty(
							"age", values.NewInt(42),
						),
						values.NewObjectProperty(
							"name", values.NewString("Jane"),
						),
						values.NewObjectProperty(
							"livesIn", values.NewObjectWith(
								values.NewObjectProperty(
									"state", values.NewString("CA"),
								),
								values.NewObjectProperty(
									"city", values.NewString("LA"),
								),
							),
						),
					),
				),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})
	})

	Convey("Merged object should be independent of source objects", t, func() {
		Convey("When array", func() {
			arr := values.NewArrayWith(values.NewInt(1), values.NewInt(2))
			obj := values.NewObjectWith(values.NewObjectProperty("arr", arr))

			actual, err := objects.MergeRecursive(context.Background(), obj)

			So(err, ShouldBeNil)
			So(actual.Compare(obj), ShouldEqual, 0)

			arr.Push(values.NewInt(0))

			So(actual.Compare(obj), ShouldNotEqual, 0)
		})

		Convey("When object", func() {
			nested := values.NewObjectWith(
				values.NewObjectProperty("nested", values.NewInt(0)),
			)
			obj := values.NewObjectWith(values.NewObjectProperty("obj", nested))

			actual, err := objects.MergeRecursive(context.Background(), obj)

			So(err, ShouldBeNil)
			So(actual.Compare(obj), ShouldEqual, 0)

			nested.Set(values.NewString("str"), values.NewInt(0))

			So(actual.Compare(obj), ShouldNotEqual, 0)
		})
	})
}
