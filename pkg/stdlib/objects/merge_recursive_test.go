package objects_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO: Fix tests
func TestMergeRecursive(t *testing.T) {
	Convey("Wrong arguments", t, func() {
		Convey("It should error when 0 arguments", func() {
			//actual, err := objects.MergeRecursive(context.Background())
			//
			//So(err, ShouldBeError)
			//So(actual.Compare(core.None), ShouldEqual, 0)
		})

		Convey("It should error when there is not object arguments", func() {
			//actual, err := objects.MergeRecursive(context.Background(), core.NewInt(0))
			//
			//So(err, ShouldBeError)
			//So(actual.Compare(core.None), ShouldEqual, 0)
			//
			//actual, err = objects.MergeRecursive(context.Background(),
			//	core.NewInt(0), runtime.NewObject(),
			//)
			//
			//So(err, ShouldBeError)
			//So(actual.Compare(core.None), ShouldEqual, 0)
		})
	})

	Convey("Merge single object", t, func() {
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("a", core.NewInt(0)),
		//)
		//expected := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("a", core.NewInt(0)),
		//)
		//
		//actual, err := objects.MergeRecursive(context.Background(), obj)
		//
		//So(err, ShouldBeNil)
		//So(actual.Compare(expected), ShouldEqual, 0)
	})

	Convey("Merge two objects", t, func() {
		Convey("When there are no common keys", func() {
			//obj1 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", core.NewInt(0)),
			//)
			//obj2 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("b", core.NewInt(1)),
			//)
			//expected := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", core.NewInt(0)),
			//	runtime.NewObjectProperty("b", core.NewInt(1)),
			//)
			//
			//actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)
			//
			//So(err, ShouldBeNil)
			//So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When objects with the same key", func() {
			//obj1 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", core.NewInt(0)),
			//	runtime.NewObjectProperty("b", core.NewInt(10)),
			//)
			//obj2 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("c", core.NewInt(1)),
			//	runtime.NewObjectProperty("b", core.NewInt(20)),
			//)
			//expected := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", core.NewInt(0)),
			//	runtime.NewObjectProperty("b", core.NewInt(20)),
			//	runtime.NewObjectProperty("c", core.NewInt(1)),
			//)
			//
			//actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)
			//
			//So(err, ShouldBeNil)
			//So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("Merge two objects with the same keys", func() {
			//obj1 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", core.NewInt(0)),
			//	runtime.NewObjectProperty("b", core.NewInt(10)),
			//)
			//obj2 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", core.NewInt(1)),
			//	runtime.NewObjectProperty("b", core.NewInt(20)),
			//)
			//expected := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", core.NewInt(1)),
			//	runtime.NewObjectProperty("b", core.NewInt(20)),
			//)
			//
			//actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)
			//
			//So(err, ShouldBeNil)
			//So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When there are nested arrays", func() {
			//obj1 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", runtime.NewArrayWith(
			//		core.NewInt(1), core.NewInt(2),
			//	)),
			//)
			//obj2 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("b", runtime.NewArrayWith(
			//		core.NewInt(1), core.NewInt(2),
			//	)),
			//)
			//expected := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("a", runtime.NewArrayWith(
			//		core.NewInt(1), core.NewInt(2),
			//	)),
			//	runtime.NewObjectProperty("b", runtime.NewArrayWith(
			//		core.NewInt(1), core.NewInt(2),
			//	)),
			//)
			//
			//actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)
			//
			//So(err, ShouldBeNil)
			//So(actual.Compare(expected), ShouldEqual, 0)
		})

		Convey("When there are nested objects (example from ArangoDB doc)", func() {
			//// { "user-1": { "name": "Jane", "livesIn": { "city": "LA" } } }
			//obj1 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty(
			//		"user-1", runtime.NewObjectWith(
			//			runtime.NewObjectProperty(
			//				"name", core.NewString("Jane"),
			//			),
			//			runtime.NewObjectProperty(
			//				"livesIn", runtime.NewObjectWith(
			//					runtime.NewObjectProperty(
			//						"city", core.NewString("LA"),
			//					),
			//				),
			//			),
			//		),
			//	),
			//)
			//// { "user-1": { "age": 42, "livesIn": { "state": "CA" } } }
			//obj2 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty(
			//		"user-1", runtime.NewObjectWith(
			//			runtime.NewObjectProperty(
			//				"age", core.NewInt(42),
			//			),
			//			runtime.NewObjectProperty(
			//				"livesIn", runtime.NewObjectWith(
			//					runtime.NewObjectProperty(
			//						"state", core.NewString("CA"),
			//					),
			//				),
			//			),
			//		),
			//	),
			//)
			//// { "user-1": { "age": 42, "livesIn": { "city": "LA", "state": "CA" }, "name": "Jane" } }
			//expected := runtime.NewObjectWith(
			//	runtime.NewObjectProperty(
			//		"user-1", runtime.NewObjectWith(
			//			runtime.NewObjectProperty(
			//				"age", core.NewInt(42),
			//			),
			//			runtime.NewObjectProperty(
			//				"name", core.NewString("Jane"),
			//			),
			//			runtime.NewObjectProperty(
			//				"livesIn", runtime.NewObjectWith(
			//					runtime.NewObjectProperty(
			//						"state", core.NewString("CA"),
			//					),
			//					runtime.NewObjectProperty(
			//						"city", core.NewString("LA"),
			//					),
			//				),
			//			),
			//		),
			//	),
			//)
			//
			//actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)
			//
			//So(err, ShouldBeNil)
			//So(actual.Compare(expected), ShouldEqual, 0)
		})
	})

	Convey("Merged object should be independent of source objects", t, func() {
		Convey("When array", func() {
			//arr := runtime.NewArrayWith(core.NewInt(1), core.NewInt(2))
			//obj := runtime.NewObjectWith(runtime.NewObjectProperty("arr", arr))
			//
			//actual, err := objects.MergeRecursive(context.Background(), obj)
			//
			//So(err, ShouldBeNil)
			//So(actual.Compare(obj), ShouldEqual, 0)
			//
			//arr.Push(core.NewInt(0))
			//
			//So(actual.Compare(obj), ShouldNotEqual, 0)
		})

		Convey("When object", func() {
			//nested := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("nested", core.NewInt(0)),
			//)
			//obj := runtime.NewObjectWith(runtime.NewObjectProperty("obj", nested))
			//
			//actual, err := objects.MergeRecursive(context.Background(), obj)
			//
			//So(err, ShouldBeNil)
			//So(actual.Compare(obj), ShouldEqual, 0)
			//
			//nested.Set(core.NewString("str"), core.NewInt(0))
			//
			//So(actual.Compare(obj), ShouldNotEqual, 0)
		})
	})
}
