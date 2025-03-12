package objects_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeRecursive(t *testing.T) {
	Convey("Wrong arguments", t, func() {
		Convey("It should error when 0 arguments", func() {
			actual, err := objects.MergeRecursive(context.Background())

			So(err, ShouldBeError)
			So(actual.Compare(core.None), ShouldEqual, 0)
		})

		Convey("It should error when there is not object arguments", func() {
			actual, err := objects.MergeRecursive(context.Background(), core.NewInt(0))

			So(err, ShouldBeError)
			So(actual.Compare(core.None), ShouldEqual, 0)

			actual, err = objects.MergeRecursive(context.Background(),
				core.NewInt(0), internal.NewObject(),
			)

			So(err, ShouldBeError)
			So(actual.Compare(core.None), ShouldEqual, 0)
		})
	})

	Convey("Merge single object", t, func() {
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("a", core.NewInt(0)),
		)
		expected := internal.NewObjectWith(
			internal.NewObjectProperty("a", core.NewInt(0)),
		)

		actual, err := objects.MergeRecursive(context.Background(), obj)

		So(err, ShouldBeNil)
		So(actual.Compare(expected), ShouldEqual, 0)
	})

	Convey("Merge two objects", t, func() {
		Convey("When there are no common keys", func() {
			obj1 := internal.NewObjectWith(
				internal.NewObjectProperty("a", core.NewInt(0)),
			)
			obj2 := internal.NewObjectWith(
				internal.NewObjectProperty("b", core.NewInt(1)),
			)
			expected := internal.NewObjectWith(
				internal.NewObjectProperty("a", core.NewInt(0)),
				internal.NewObjectProperty("b", core.NewInt(1)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When objects with the same key", func() {
			obj1 := internal.NewObjectWith(
				internal.NewObjectProperty("a", core.NewInt(0)),
				internal.NewObjectProperty("b", core.NewInt(10)),
			)
			obj2 := internal.NewObjectWith(
				internal.NewObjectProperty("c", core.NewInt(1)),
				internal.NewObjectProperty("b", core.NewInt(20)),
			)
			expected := internal.NewObjectWith(
				internal.NewObjectProperty("a", core.NewInt(0)),
				internal.NewObjectProperty("b", core.NewInt(20)),
				internal.NewObjectProperty("c", core.NewInt(1)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("Merge two objects with the same keys", func() {
			obj1 := internal.NewObjectWith(
				internal.NewObjectProperty("a", core.NewInt(0)),
				internal.NewObjectProperty("b", core.NewInt(10)),
			)
			obj2 := internal.NewObjectWith(
				internal.NewObjectProperty("a", core.NewInt(1)),
				internal.NewObjectProperty("b", core.NewInt(20)),
			)
			expected := internal.NewObjectWith(
				internal.NewObjectProperty("a", core.NewInt(1)),
				internal.NewObjectProperty("b", core.NewInt(20)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When there are nested arrays", func() {
			obj1 := internal.NewObjectWith(
				internal.NewObjectProperty("a", internal.NewArrayWith(
					core.NewInt(1), core.NewInt(2),
				)),
			)
			obj2 := internal.NewObjectWith(
				internal.NewObjectProperty("b", internal.NewArrayWith(
					core.NewInt(1), core.NewInt(2),
				)),
			)
			expected := internal.NewObjectWith(
				internal.NewObjectProperty("a", internal.NewArrayWith(
					core.NewInt(1), core.NewInt(2),
				)),
				internal.NewObjectProperty("b", internal.NewArrayWith(
					core.NewInt(1), core.NewInt(2),
				)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When there are nested objects (example from ArangoDB doc)", func() {
			// { "user-1": { "name": "Jane", "livesIn": { "city": "LA" } } }
			obj1 := internal.NewObjectWith(
				internal.NewObjectProperty(
					"user-1", internal.NewObjectWith(
						internal.NewObjectProperty(
							"name", core.NewString("Jane"),
						),
						internal.NewObjectProperty(
							"livesIn", internal.NewObjectWith(
								internal.NewObjectProperty(
									"city", core.NewString("LA"),
								),
							),
						),
					),
				),
			)
			// { "user-1": { "age": 42, "livesIn": { "state": "CA" } } }
			obj2 := internal.NewObjectWith(
				internal.NewObjectProperty(
					"user-1", internal.NewObjectWith(
						internal.NewObjectProperty(
							"age", core.NewInt(42),
						),
						internal.NewObjectProperty(
							"livesIn", internal.NewObjectWith(
								internal.NewObjectProperty(
									"state", core.NewString("CA"),
								),
							),
						),
					),
				),
			)
			// { "user-1": { "age": 42, "livesIn": { "city": "LA", "state": "CA" }, "name": "Jane" } }
			expected := internal.NewObjectWith(
				internal.NewObjectProperty(
					"user-1", internal.NewObjectWith(
						internal.NewObjectProperty(
							"age", core.NewInt(42),
						),
						internal.NewObjectProperty(
							"name", core.NewString("Jane"),
						),
						internal.NewObjectProperty(
							"livesIn", internal.NewObjectWith(
								internal.NewObjectProperty(
									"state", core.NewString("CA"),
								),
								internal.NewObjectProperty(
									"city", core.NewString("LA"),
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
			arr := internal.NewArrayWith(core.NewInt(1), core.NewInt(2))
			obj := internal.NewObjectWith(internal.NewObjectProperty("arr", arr))

			actual, err := objects.MergeRecursive(context.Background(), obj)

			So(err, ShouldBeNil)
			So(actual.Compare(obj), ShouldEqual, 0)

			arr.Push(core.NewInt(0))

			So(actual.Compare(obj), ShouldNotEqual, 0)
		})

		Convey("When object", func() {
			nested := internal.NewObjectWith(
				internal.NewObjectProperty("nested", core.NewInt(0)),
			)
			obj := internal.NewObjectWith(internal.NewObjectProperty("obj", nested))

			actual, err := objects.MergeRecursive(context.Background(), obj)

			So(err, ShouldBeNil)
			So(actual.Compare(obj), ShouldEqual, 0)

			nested.Set(core.NewString("str"), core.NewInt(0))

			So(actual.Compare(obj), ShouldNotEqual, 0)
		})
	})
}
