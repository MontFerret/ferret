package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeys(t *testing.T) {
	Convey("Keys(obj, false) should return 'a', 'c', 'b' in any order", t, func() {
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("a", core.NewInt(0)),
			internal.NewObjectProperty("b", core.NewInt(1)),
			internal.NewObjectProperty("c", core.NewInt(2)),
		)

		keys, err := objects.Keys(context.Background(), obj)
		keysArray := keys.(*internal.Array)

		So(err, ShouldEqual, nil)
		So(keysArray.Type().Equals(types.Array), ShouldBeTrue)
		So(keysArray.Length(), ShouldEqual, 3)

		for _, k := range []string{"b", "a", "c"} {
			iof := keysArray.IndexOf(core.NewString(k))
			So(iof, ShouldNotEqual, -1)
		}
	})

	Convey("Keys(obj, false) should return ['a', 'b', 'c']", t, func() {
		obj := internal.NewObjectWith(
			internal.NewObjectProperty("b", core.NewInt(0)),
			internal.NewObjectProperty("a", core.NewInt(1)),
			internal.NewObjectProperty("c", core.NewInt(3)),
		)

		keys, err := objects.Keys(context.Background(), obj, core.NewBoolean(true))
		keysArray := keys.(*internal.Array)

		So(err, ShouldEqual, nil)

		for idx, key := range []string{"a", "b", "c"} {
			So(keysArray.Get(core.NewInt(idx)), ShouldEqual, core.NewString(key))
		}
	})

	Convey("When there are no keys", t, func() {
		obj := internal.NewObject()

		keys, err := objects.Keys(context.Background(), obj, core.NewBoolean(true))
		keysArray := keys.(*internal.Array)

		So(err, ShouldEqual, nil)
		So(keysArray.Length(), ShouldEqual, core.NewInt(0))
		So(int(keysArray.Length()), ShouldEqual, 0)

		keys, err = objects.Keys(context.Background(), obj, core.NewBoolean(false))
		keysArray = keys.(*internal.Array)

		So(err, ShouldEqual, nil)
		So(keysArray.Length(), ShouldEqual, core.NewInt(0))
		So(int(keysArray.Length()), ShouldEqual, 0)
	})

	Convey("When not enough arguments", t, func() {
		_, err := objects.Keys(context.Background())

		So(err, ShouldBeError)
	})

	Convey("When first argument isn't object", t, func() {
		notObj := core.NewInt(0)

		_, err := objects.Keys(context.Background(), notObj)

		So(err, ShouldBeError)
	})

	Convey("When second argument isn't boolean", t, func() {
		obj := internal.NewObject()

		_, err := objects.Keys(context.Background(), obj, obj)

		So(err, ShouldBeError)
	})

}
