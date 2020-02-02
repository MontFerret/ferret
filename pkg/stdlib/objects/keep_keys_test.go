package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKeepKeys(t *testing.T) {
	Convey("When not enough arguments)", t, func() {
		// there is no object
		obj, err := objects.KeepKeys(context.Background())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)

		// there are no keys
		obj, err = objects.KeepKeys(context.Background(), values.NewObject())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("When first argument isn't object", t, func() {
		obj, err := objects.KeepKeys(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("When wrong keys arguments", t, func() {
		obj, err := objects.KeepKeys(context.Background(), values.NewObject(), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)

		// looks like a valid case
		// but there is another argument besides an array
		obj, err = objects.KeepKeys(context.Background(), values.NewObject(), values.NewArray(0), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("Result object is independent of the source object", t, func() {
		arr := values.NewArrayWith(values.Int(0))
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", arr),
		)
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewArrayWith(values.Int(0))),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, values.NewString("a"))

		So(err, ShouldBeNil)

		arr.Push(values.NewInt(1))

		So(afterKeepKeys.Compare(resultObj), ShouldEqual, 0)
	})
}

func TestKeepKeysStrings(t *testing.T) {
	Convey("KeepKeys key 'a'", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, values.NewString("a"))

		So(err, ShouldEqual, nil)
		So(afterKeepKeys.Compare(resultObj), ShouldEqual, 0)
	})

	Convey("KeepKeys key doesn't exists", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		resultObj := values.NewObject()

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, values.NewString("c"))

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*values.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys when there are more keys than object properties", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj,
			values.NewString("a"), values.NewString("b"), values.NewString("c"),
		)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*values.Object), resultObj), ShouldEqual, true)
	})
}

func TestKeepKeysArray(t *testing.T) {
	Convey("KeepKeys array", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		keys := values.NewArrayWith(values.NewString("a"))
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*values.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys empty array", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		keys := values.NewArray(0)
		resultObj := values.NewObject()

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*values.Object), resultObj), ShouldEqual, true)
	})

	Convey("KeepKeys when there are more keys than object properties", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		keys := values.NewArrayWith(
			values.NewString("a"), values.NewString("b"), values.NewString("c"),
		)
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeepKeys.(*values.Object), resultObj), ShouldEqual, true)
	})

	Convey("When there is not string key", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		keys := values.NewArrayWith(
			values.NewString("a"),
			values.NewInt(0),
		)

		afterKeepKeys, err := objects.KeepKeys(context.Background(), obj, keys)

		So(err, ShouldBeError)
		So(afterKeepKeys, ShouldEqual, values.None)
	})
}

func isEqualObjects(obj1 *values.Object, obj2 *values.Object) bool {
	var val1 core.Value
	var val2 core.Value

	for _, key := range obj1.Keys() {
		val1, _ = obj1.Get(key)
		val2, _ = obj2.Get(key)
		if val1.Compare(val2) != 0 {
			return false
		}
	}
	for _, key := range obj2.Keys() {
		val1, _ = obj1.Get(key)
		val2, _ = obj2.Get(key)
		if val2.Compare(val1) != 0 {
			return false
		}
	}
	return true
}
