package objects_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeys(t *testing.T) {
	Convey("Keys(obj, false) should return 'a', 'c', 'b' in any order", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("a", values.NewInt(0)),
			values.NewObjectProperty("b", values.NewInt(1)),
			values.NewObjectProperty("c", values.NewInt(2)),
		)

		keys, err := objects.Keys(context.Background(), obj)
		keysArray := keys.(*values.Array)

		So(err, ShouldEqual, nil)
		So(keysArray.Type().Equals(types.Array), ShouldBeTrue)
		So(keysArray.Length(), ShouldEqual, 3)

		for _, k := range []string{"b", "a", "c"} {
			iof := keysArray.IndexOf(values.NewString(k))
			So(iof, ShouldNotEqual, -1)
		}
	})

	Convey("Keys(obj, false) should return ['a', 'b', 'c']", t, func() {
		obj := values.NewObjectWith(
			values.NewObjectProperty("b", values.NewInt(0)),
			values.NewObjectProperty("a", values.NewInt(1)),
			values.NewObjectProperty("c", values.NewInt(3)),
		)

		keys, err := objects.Keys(context.Background(), obj, values.NewBoolean(true))
		keysArray := keys.(*values.Array)

		So(err, ShouldEqual, nil)

		for idx, key := range []string{"a", "b", "c"} {
			So(keysArray.Get(values.NewInt(idx)), ShouldEqual, values.NewString(key))
		}
	})

	Convey("When there are no keys", t, func() {
		obj := values.NewObject()

		keys, err := objects.Keys(context.Background(), obj, values.NewBoolean(true))
		keysArray := keys.(*values.Array)

		So(err, ShouldEqual, nil)
		So(keysArray.Length(), ShouldEqual, values.NewInt(0))
		So(int(keysArray.Length()), ShouldEqual, 0)

		keys, err = objects.Keys(context.Background(), obj, values.NewBoolean(false))
		keysArray = keys.(*values.Array)

		So(err, ShouldEqual, nil)
		So(keysArray.Length(), ShouldEqual, values.NewInt(0))
		So(int(keysArray.Length()), ShouldEqual, 0)
	})

	Convey("When not enough arguments", t, func() {
		_, err := objects.Keys(context.Background())

		So(err, ShouldBeError)
	})

	Convey("When first argument isn't object", t, func() {
		notObj := values.NewInt(0)

		_, err := objects.Keys(context.Background(), notObj)

		So(err, ShouldBeError)
	})

	Convey("When second argument isn't boolean", t, func() {
		obj := values.NewObject()

		_, err := objects.Keys(context.Background(), obj, obj)

		So(err, ShouldBeError)
	})

}
