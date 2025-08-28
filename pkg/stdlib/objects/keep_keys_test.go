package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"
)

// TODO: Fix tests
func TestKeepKeys(t *testing.T) {
	Convey("When not enough arguments)", t, func() {
		// there is no object
		obj, err := objects.KeepKeys(context.Background())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, runtime.None)

		// there are no keys
		obj, err = objects.KeepKeys(context.Background(), runtime.NewObject())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, runtime.None)
	})

	Convey("When first argument isn't object", t, func() {
		obj, err := objects.KeepKeys(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, runtime.None)
	})

	Convey("When wrong keys arguments", t, func() {
		obj, err := objects.KeepKeys(context.Background(), runtime.NewObject(), runtime.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, runtime.None)

		// looks like a valid case
		// but there is another argument besides an array
		obj, err = objects.KeepKeys(context.Background(), runtime.NewObject(), runtime.NewArray(0), runtime.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, runtime.None)
	})

	Convey("Result object is independent of the source object", t, func() {
		arr := runtime.NewArrayWith(runtime.Int(0))
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", arr),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, runtime.NewString("a"))

		So(err, ShouldBeNil)

		arr.Add(context.Background(), runtime.NewInt(1))

		resultObj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewArrayWith(runtime.Int(0))),
		)
		So(runtime.CompareValues(afterKeepKeys, resultObj), ShouldEqual, 0)
	})
}

func TestKeepKeysStrings(t *testing.T) {
	Convey("KeepKeys key 'a'", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)
		resultObj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, runtime.NewString("a"))

		So(err, ShouldEqual, nil)
		So(runtime.CompareValues(afterKeepKeys, resultObj), ShouldEqual, 0)
	})

	Convey("KeepKeys key doesn't exists", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)
		resultObj := runtime.NewObject()

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, runtime.NewString("c"))

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys when there are more keys than object properties", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)
		resultObj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj,
			runtime.NewString("a"), runtime.NewString("b"), runtime.NewString("c"),
		)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})
}

func TestKeepKeysArray(t *testing.T) {
	Convey("KeepKeys array", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)
		keys := runtime.NewArrayWith(runtime.NewString("a"))
		resultObj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys empty array", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)
		keys := runtime.NewArray(0)
		resultObj := runtime.NewObject()

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys when there are more keys than object properties", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)
		keys := runtime.NewArrayWith(
			runtime.NewString("a"), runtime.NewString("b"), runtime.NewString("c"),
		)
		resultObj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})

	Convey("When there is not string key", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("b", runtime.NewString("string")),
		)
		keys := runtime.NewArrayWith(
			runtime.NewString("a"),
			runtime.NewInt(0),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldBeError)
		So(afterKeepKeys, ShouldEqual, runtime.None)
	})
}

func isEqualObjects(obj1 *runtime.Object, obj2 *runtime.Object) bool {
	// Use the built-in Compare method
	return runtime.CompareValues(obj1, obj2) == 0
}
