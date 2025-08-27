package objects_test

import (
	"testing"
)

// TODO: Fix the tests
func TestZip(t *testing.T) {
	//Convey("Invalid arguments", t, func() {
	//	Convey("When there are no arguments", func() {
	//		actual, err := objects.Zip(context.Background())
	//		expected := runtime.None
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//	})
	//
	//	Convey("When single argument", func() {
	//		actual, err := objects.Zip(context.Background(), runtime.NewArray(0))
	//		expected := runtime.None
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//
	//		actual, err = objects.Zip(context.Background(), runtime.NewInt(0))
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//	})
	//
	//	Convey("When too many arguments", func() {
	//		actual, err := objects.Zip(context.Background(),
	//			runtime.NewArray(0), runtime.NewArray(0), runtime.NewArray(0))
	//		expected := runtime.None
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//	})
	//
	//	Convey("When there is not array argument", func() {
	//		actual, err := objects.Zip(context.Background(), runtime.NewArray(0), runtime.NewInt(0))
	//		expected := runtime.None
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//
	//		actual, err = objects.Zip(context.Background(), runtime.NewInt(0), runtime.NewArray(0))
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//	})
	//
	//	Convey("When there is not string element into keys array", func() {
	//		keys := runtime.NewArrayWith(runtime.NewInt(0))
	//		vals := runtime.NewArrayWith(runtime.NewString("v1"))
	//		expected := runtime.None
	//
	//		actual, err := objects.Zip(context.Background(), keys, vals)
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//	})
	//
	//	Convey("When 1 key and 0 values", func() {
	//		keys := runtime.NewArrayWith(runtime.NewString("k1"))
	//		vals := runtime.NewArray(0)
	//		expected := runtime.None
	//
	//		actual, err := objects.Zip(context.Background(), keys, vals)
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//	})
	//
	//	Convey("When 0 keys and 1 values", func() {
	//		keys := runtime.NewArray(0)
	//		vals := runtime.NewArrayWith(runtime.NewString("v1"))
	//		expected := runtime.None
	//
	//		actual, err := objects.Zip(context.Background(), keys, vals)
	//
	//		So(err, ShouldBeError)
	//		So(actual.Compare(expected), ShouldEqual, 0)
	//	})
	//})
	//
	//Convey("Zip 2 keys and 2 values", t, func() {
	//	keys := runtime.NewArrayWith(
	//		runtime.NewString("k1"),
	//		runtime.NewString("k2"),
	//	)
	//	vals := runtime.NewArrayWith(
	//		runtime.NewString("v1"),
	//		runtime.NewInt(2),
	//	)
	//	expected := runtime.NewObjectWith(
	//		runtime.NewObjectProperty("k1", runtime.NewString("v1")),
	//		runtime.NewObjectProperty("k2", runtime.NewInt(2)),
	//	)
	//
	//	actual, err := objects.Zip(context.Background(), keys, vals)
	//
	//	So(err, ShouldBeNil)
	//	So(actual.Compare(expected), ShouldEqual, 0)
	//})
	//
	//Convey("Zip 3 keys and 3 values. 1 key repeats", t, func() {
	//	keys := runtime.NewArrayWith(
	//		runtime.NewString("k1"),
	//		runtime.NewString("k2"),
	//		runtime.NewString("k1"),
	//	)
	//	vals := runtime.NewArrayWith(
	//		runtime.NewInt(1),
	//		runtime.NewInt(2),
	//		runtime.NewInt(3),
	//	)
	//	expected := runtime.NewObjectWith(
	//		runtime.NewObjectProperty("k1", runtime.NewInt(1)),
	//		runtime.NewObjectProperty("k2", runtime.NewInt(2)),
	//	)
	//
	//	actual, err := objects.Zip(context.Background(), keys, vals)
	//
	//	So(err, ShouldBeNil)
	//	So(actual.Compare(expected), ShouldEqual, 0)
	//})
}
