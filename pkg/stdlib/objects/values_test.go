package objects_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO: Fix the tests
func TestValues(t *testing.T) {
	Convey("Invalid arguments", t, func() {
		Convey("When there is no arguments", func() {
			//actual, err := objects.Values(context.Background())
			//
			//So(err, ShouldBeError)
			//So(actual.Compare(core.None), ShouldEqual, 0)
		})

		Convey("When 2 arguments", func() {
			//obj := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("k1", core.NewInt(0)),
			//	runtime.NewObjectProperty("k2", core.NewInt(1)),
			//)
			//
			//actual, err := objects.Values(context.Background(), obj, obj)
			//
			//So(err, ShouldBeError)
			//So(actual.Compare(core.None), ShouldEqual, 0)
			//
			//actual, err = objects.Values(context.Background(), obj, core.NewInt(0))
			//
			//So(err, ShouldBeError)
			//So(actual.Compare(core.None), ShouldEqual, 0)
		})

		Convey("When there is not object argument", func() {
			//actual, err := objects.Values(context.Background(), core.NewInt(0))
			//
			//So(err, ShouldBeError)
			//So(actual.Compare(core.None), ShouldEqual, 0)
		})
	})

	Convey("When simple type attributes (same type)", t, func() {
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("k1", core.NewInt(0)),
		//	runtime.NewObjectProperty("k2", core.NewInt(1)),
		//)
		//expected := runtime.NewArrayWith(
		//	core.NewInt(0), core.NewInt(1),
		//).Sort()
		//
		//actual, err := objects.Values(context.Background(), obj)
		//actualSorted := actual.(*runtime.Array).Sort()
		//
		//So(err, ShouldBeNil)
		//So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When simple type attributes (different types)", t, func() {
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("k1", core.NewInt(0)),
		//	runtime.NewObjectProperty("k2", core.NewString("v2")),
		//)
		//expected := runtime.NewArrayWith(
		//	core.NewInt(0), core.NewString("v2"),
		//).Sort()
		//
		//actual, err := objects.Values(context.Background(), obj)
		//actualSorted := actual.(*runtime.Array).Sort()
		//
		//So(err, ShouldBeNil)
		//So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (array)", t, func() {
		//arr1 := runtime.NewArrayWith(
		//	core.NewInt(0), core.NewInt(1),
		//)
		//arr2 := runtime.NewArrayWith(
		//	core.NewInt(2), core.NewInt(3),
		//)
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("k1", arr1),
		//	runtime.NewObjectProperty("k2", arr2),
		//)
		//expected := runtime.NewArrayWith(arr1, arr2).Sort()
		//
		//actual, err := objects.Values(context.Background(), obj)
		//actualSorted := actual.(*runtime.Array).Sort()
		//
		//So(err, ShouldBeNil)
		//So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (object)", t, func() {
		//obj1 := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("int0", core.NewInt(0)),
		//)
		//obj2 := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("int1", core.NewInt(1)),
		//)
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("k1", obj1),
		//	runtime.NewObjectProperty("k2", obj2),
		//)
		//expected := runtime.NewArrayWith(obj1, obj2).Sort()
		//
		//actual, err := objects.Values(context.Background(), obj)
		//actualSorted := actual.(*runtime.Array).Sort()
		//
		//So(err, ShouldBeNil)
		//So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When complex type attributes (object and array)", t, func() {
		//obj1 := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("k1", core.NewInt(0)),
		//)
		//arr1 := runtime.NewArrayWith(
		//	core.NewInt(0), core.NewInt(1),
		//)
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("obj", obj1),
		//	runtime.NewObjectProperty("arr", arr1),
		//)
		//expected := runtime.NewArrayWith(obj1, arr1).Sort()
		//
		//actual, err := objects.Values(context.Background(), obj)
		//actualSorted := actual.(*runtime.Array).Sort()
		//
		//So(err, ShouldBeNil)
		//So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("When both type attributes", t, func() {
		//obj1 := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("k1", core.NewInt(0)),
		//)
		//arr1 := runtime.NewArrayWith(
		//	core.NewInt(0), core.NewInt(1),
		//)
		//int1 := core.NewInt(0)
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("obj", obj1),
		//	runtime.NewObjectProperty("arr", arr1),
		//	runtime.NewObjectProperty("int", int1),
		//)
		//expected := runtime.NewArrayWith(obj1, arr1, int1).Sort()
		//
		//actual, err := objects.Values(context.Background(), obj)
		//actualSorted := actual.(*runtime.Array).Sort()
		//
		//So(err, ShouldBeNil)
		//So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("Result is independent on the source object (array)", t, func() {
		//arr := runtime.NewArrayWith(core.NewInt(0))
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("arr", arr),
		//)
		//expected := runtime.NewArrayWith(
		//	runtime.NewArrayWith(
		//		core.NewInt(0),
		//	),
		//)
		//
		//actual, err := objects.Values(context.Background(), obj)
		//actualSorted := actual.(*runtime.Array).Sort()
		//
		//So(err, ShouldBeNil)
		//
		//arr.Push(core.NewInt(1))
		//
		//So(actualSorted.Compare(expected), ShouldEqual, 0)
	})

	Convey("Result is independent on the source object (object)", t, func() {
		//nested := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("int", core.NewInt(0)),
		//)
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("nested", nested),
		//)
		//expected := runtime.NewArrayWith(
		//	runtime.NewObjectWith(
		//		runtime.NewObjectProperty("int", core.NewInt(0)),
		//	),
		//)
		//
		//actual, err := objects.Values(context.Background(), obj)
		//actualSorted := actual.(*runtime.Array).Sort()
		//
		//So(err, ShouldBeNil)
		//
		//nested.Set("new", core.NewInt(1))
		//
		//So(actualSorted.Compare(expected), ShouldEqual, 0)
	})
}

func TestValuesStress(t *testing.T) {
	Convey("Stress", t, func() {
		for i := 0; i < 100; i++ {
			//obj1 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("int0", core.NewInt(0)),
			//)
			//obj2 := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("int1", core.NewInt(1)),
			//)
			//obj := runtime.NewObjectWith(
			//	runtime.NewObjectProperty("k1", obj1),
			//	runtime.NewObjectProperty("k2", obj2),
			//)
			//expected := runtime.NewArrayWith(obj2, obj1).Sort()
			//
			//actual, err := objects.Values(context.Background(), obj)
			//actualSorted := actual.(*runtime.Array).Sort()
			//
			//So(err, ShouldBeNil)
			//So(actualSorted.Length(), ShouldEqual, expected.Length())
			//So(actualSorted.Compare(expected), ShouldEqual, 0)
		}
	})
}
