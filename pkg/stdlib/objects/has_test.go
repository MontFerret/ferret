package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHas(t *testing.T) {
	Convey("When key exists", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("key", values.NewString("val")),
		)

		val, err := objects.Has(context.Background(), obj, values.NewString("key"))
		valBool := val.(values.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, values.NewBoolean(true))
		So(bool(valBool), ShouldEqual, true)
	})

	Convey("When key doesn't exists", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("anyOtherKey", values.NewString("val")),
		)

		val, err := objects.Has(context.Background(), obj, values.NewString("key"))
		valBool := val.(values.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, values.NewBoolean(false))
		So(bool(valBool), ShouldEqual, false)
	})

	Convey("When there are no keys", t, func() {
		obj := values.NewObject()

		val, err := objects.Has(context.Background(), obj, values.NewString("key"))
		valBool := val.(values.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, values.NewBoolean(false))
		So(bool(valBool), ShouldEqual, false)
	})

	Convey("Not enought arguments", t, func() {
		val, err := objects.Has(context.Background())

		So(err, ShouldBeError)
		So(val, ShouldEqual, values.None)

		val, err = objects.Has(context.Background(), values.NewObject())

		So(err, ShouldBeError)
		So(val, ShouldEqual, values.None)
	})

	Convey("When keyName isn't string", t, func() {
		obj := values.NewObject()
		key := values.NewInt(1)

		val, err := objects.Has(context.Background(), obj, key)

		So(err, ShouldBeError)
		So(val, ShouldEqual, values.None)
	})

	Convey("When first argument isn't object", t, func() {
		notObj := values.NewInt(1)

		val, err := objects.Has(context.Background(), notObj, values.NewString("key"))

		So(err, ShouldBeError)
		So(val, ShouldEqual, values.None)
	})
}
