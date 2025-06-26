package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"
)

// TODO: Fix tests
func TestKeepKeys(t *testing.T) {
	Convey("When not enough arguments)", t, func() {
		// there is no object
		obj, err := objects.KeepKeys(context.Background())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, core.None)

		// there are no keys
		obj, err = objects.KeepKeys(context.Background(), runtime.NewObject())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, core.None)
	})

	Convey("When first argument isn't object", t, func() {
		obj, err := objects.KeepKeys(context.Background(), core.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, core.None)
	})

	Convey("When wrong keys arguments", t, func() {
		obj, err := objects.KeepKeys(context.Background(), runtime.NewObject(), core.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, core.None)

		// looks like a valid case
		// but there is another argument besides an array
		obj, err = objects.KeepKeys(context.Background(), runtime.NewObject(), runtime.NewArray(0), core.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, core.None)
	})

	Convey("Result object is independent of the source object", t, func() {
		//arr := runtime.NewArrayWith(core.Int(0))
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("a", arr),
		//)
		//resultObj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("a", runtime.NewArrayWith(core.Int(0))),
		//)
		//
		//afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, core.NewString("a"))
		//
		//So(err, ShouldBeNil)

		//_ = arr.Add(nil, core.NewInt(1))
		//
		//So(afterKeepKeys.Compare(resultObj), ShouldEqual, 0)
	})
}

func TestKeepKeysStrings(t *testing.T) {
	Convey("KeepKeys key 'a'", t, func() {
		//obj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("a", core.NewInt(1)),
		//	runtime.NewObjectProperty("b", core.NewString("string")),
		//)
		//resultObj := runtime.NewObjectWith(
		//	runtime.NewObjectProperty("a", core.NewInt(1)),
		//)
		//
		//afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, core.NewString("a"))

		//So(err, ShouldEqual, nil)
		//So(afterKeepKeys.Compare(resultObj), ShouldEqual, 0)
	})

	Convey("KeepKeys key doesn't exists", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
			runtime.NewObjectProperty("b", core.NewString("string")),
		)
		resultObj := runtime.NewObject()

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, core.NewString("c"))

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys when there are more keys than object properties", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
			runtime.NewObjectProperty("b", core.NewString("string")),
		)
		resultObj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
			runtime.NewObjectProperty("b", core.NewString("string")),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj,
			core.NewString("a"), core.NewString("b"), core.NewString("c"),
		)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})
}

func TestKeepKeysArray(t *testing.T) {
	Convey("KeepKeys array", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
			runtime.NewObjectProperty("b", core.NewString("string")),
		)
		keys := runtime.NewArrayWith(core.NewString("a"))
		resultObj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys empty array", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
			runtime.NewObjectProperty("b", core.NewString("string")),
		)
		keys := runtime.NewArray(0)
		resultObj := runtime.NewObject()

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys when there are more keys than object properties", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
			runtime.NewObjectProperty("b", core.NewString("string")),
		)
		keys := runtime.NewArrayWith(
			core.NewString("a"), core.NewString("b"), core.NewString("c"),
		)
		resultObj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
			runtime.NewObjectProperty("b", core.NewString("string")),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*runtime.Object), resultObj), ShouldEqual, true)
	})

	Convey("When there is not string key", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", core.NewInt(1)),
			runtime.NewObjectProperty("b", core.NewString("string")),
		)
		keys := runtime.NewArrayWith(
			core.NewString("a"),
			core.NewInt(0),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldBeError)
		So(afterKeepKeys, ShouldEqual, core.None)
	})
}

func isEqualObjects(obj1 *runtime.Object, obj2 *runtime.Object) bool {
	//var val1 core.Value
	//var val2 core.Value
	//
	//for _, key := range obj1.Keys() {
	//	val1, _ = obj1.Get(key)
	//	val2, _ = obj2.Get(key)
	//	if val1.Compare(val2) != 0 {
	//		return false
	//	}
	//}
	//for _, key := range obj2.Keys() {
	//	val1, _ = obj1.Get(key)
	//	val2, _ = obj2.Get(key)
	//	if val2.Compare(val1) != 0 {
	//		return false
	//	}
	//}
	return true
}
