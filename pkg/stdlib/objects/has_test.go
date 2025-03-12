package objects_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHas(t *testing.T) {
	Convey("When key exists", t, func() {
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("key", core.NewString("val")),
		)

		val, err := objects.Has(context.Background(), obj, core.NewString("key"))
		valBool := val.(core.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, core.NewBoolean(true))
		So(bool(valBool), ShouldEqual, true)
	})

	Convey("When key doesn't exists", t, func() {
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("anyOtherKey", core.NewString("val")),
		)

		val, err := objects.Has(context.Background(), obj, core.NewString("key"))
		valBool := val.(core.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, core.NewBoolean(false))
		So(bool(valBool), ShouldEqual, false)
	})

	Convey("When there are no keys", t, func() {
		obj := internal.NewObject()

		val, err := objects.Has(context.Background(), obj, core.NewString("key"))
		valBool := val.(core.Boolean)

		So(err, ShouldEqual, nil)
		So(valBool, ShouldEqual, core.NewBoolean(false))
		So(bool(valBool), ShouldEqual, false)
	})

	Convey("Not enought arguments", t, func() {
		val, err := objects.Has(context.Background())

		So(err, ShouldBeError)
		So(val, ShouldEqual, core.None)

		val, err = objects.Has(context.Background(), internal.NewObject())

		So(err, ShouldBeError)
		So(val, ShouldEqual, core.None)
	})

	Convey("When keyName isn't string", t, func() {
		obj := internal.NewObject()
		key := core.NewInt(1)

		val, err := objects.Has(context.Background(), obj, key)

		So(err, ShouldBeError)
		So(val, ShouldEqual, core.None)
	})

	Convey("When first argument isn't object", t, func() {
		notObj := core.NewInt(1)

		val, err := objects.Has(context.Background(), notObj, core.NewString("key"))

		So(err, ShouldBeError)
		So(val, ShouldEqual, core.None)
	})
}
