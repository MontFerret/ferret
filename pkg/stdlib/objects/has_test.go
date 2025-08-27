package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHas(t *testing.T) {
	Convey("When key exists", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("key", runtime.NewString("val")),
		)

		val, err := objects.Has(context.Background(), obj, runtime.NewString("key"))
		valBool := val.(runtime.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, runtime.NewBoolean(true))
		So(bool(valBool), ShouldEqual, true)
	})

	Convey("When key doesn't exists", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("anyOtherKey", runtime.NewString("val")),
		)

		val, err := objects.Has(context.Background(), obj, runtime.NewString("key"))
		valBool := val.(runtime.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, runtime.NewBoolean(false))
		So(bool(valBool), ShouldEqual, false)
	})

	Convey("When there are no keys", t, func() {
		obj := runtime.NewObject()

		val, err := objects.Has(context.Background(), obj, runtime.NewString("key"))
		valBool := val.(runtime.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, runtime.NewBoolean(false))
		So(bool(valBool), ShouldEqual, false)
	})

	Convey("Not enought arguments", t, func() {
		val, err := objects.Has(context.Background())

		So(err, ShouldBeError)
		So(val, ShouldEqual, runtime.None)

		val, err = objects.Has(context.Background(), runtime.NewObject())

		So(err, ShouldBeError)
		So(val, ShouldEqual, runtime.None)
	})

	Convey("When keyName isn't string", t, func() {
		obj := runtime.NewObject()
		key := runtime.NewInt(1)

		val, err := objects.Has(context.Background(), obj, key)

		So(err, ShouldBeError)
		So(val, ShouldEqual, runtime.None)
	})

	Convey("When first argument isn't object", t, func() {
		notObj := runtime.NewInt(1)

		val, err := objects.Has(context.Background(), notObj, runtime.NewString("key"))

		So(err, ShouldBeError)
		So(val, ShouldEqual, runtime.None)
	})
}
