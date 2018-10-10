package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKeep(t *testing.T) {
	Convey("When not enought arguments)", t, func() {
		// there is no object
		obj, err := objects.Keep(context.Background())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)

		// there are no keys
		obj, err = objects.Keep(context.Background(), values.NewObject())

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("When first argument isn't object", t, func() {
		obj, err := objects.Keep(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})

	Convey("When wrong keys arguments", t, func() {
		obj, err := objects.Keep(context.Background(), values.NewObject(), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)

		// looks like a valid case
		// but there is another argument besides an array
		obj, err = objects.Keep(context.Background(), values.NewObject(), values.NewArray(0), values.NewInt(0))

		So(err, ShouldBeError)
		So(obj, ShouldEqual, values.None)
	})
}

func TestKeepStrings(t *testing.T) {
	Convey("Keep key 'a'", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
		)

		afterKeep, err := objects.Keep(context.Background(), obj, values.NewString("a"))

		So(err, ShouldEqual, nil)
		So(afterKeep.Compare(resultObj), ShouldEqual, 0)
	})

	Convey("Keep key doesn't exists", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		resultObj := values.NewObject()

		afterKeep, err := objects.Keep(context.Background(), obj, values.NewString("c"))

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeep.(*values.Object), resultObj), ShouldEqual, true)
	})

	Convey("Keep when there are more keys than object properties", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)

		afterKeep, err := objects.Keep(context.Background(), obj,
			values.NewString("a"), values.NewString("b"), values.NewString("c"),
		)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeep.(*values.Object), resultObj), ShouldEqual, true)
	})
}

func TestKeepArray(t *testing.T) {
	Convey("Keep array", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		keys := values.NewArrayWith(values.NewString("a"))
		resultObj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
		)

		afterKeep, err := objects.Keep(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeep.(*values.Object), resultObj), ShouldEqual, true)
	})

	Convey("Keep empty array", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("b", values.NewString("string")),
		)
		keys := values.NewArray(0)
		resultObj := values.NewObject()

		afterKeep, err := objects.Keep(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeep.(*values.Object), resultObj), ShouldEqual, true)
	})

	Convey("Keep when there are more keys than object properties", t, func() {
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

		afterKeep, err := objects.Keep(context.Background(), obj, keys)

		So(err, ShouldEqual, nil)
		So(isEqualObjects(afterKeep.(*values.Object), resultObj), ShouldEqual, true)
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

		afterKeep, err := objects.Keep(context.Background(), obj, keys)

		So(err, ShouldBeError)
		So(afterKeep, ShouldEqual, values.None)
	})
}

func isEqualObjects(obj1 *values.Object, obj2 *values.Object) bool {
	var val1 core.Value
	var val2 core.Value

	for _, key := range obj1.Keys() {
		val1, _ = obj1.Get(values.NewString(key))
		val2, _ = obj2.Get(values.NewString(key))
		if val1.Compare(val2) != 0 {
			return false
		}
	}
	for _, key := range obj2.Keys() {
		val1, _ = obj1.Get(values.NewString(key))
		val2, _ = obj2.Get(values.NewString(key))
		if val2.Compare(val1) != 0 {
			return false
		}
	}
	return true
}
