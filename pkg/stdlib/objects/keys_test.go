package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeys(t *testing.T) {
	Convey("Keys(obj, false) should return 'a', 'c', 'b' in any order", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(0)),
			runtime.NewObjectProperty("b", runtime.NewInt(1)),
			runtime.NewObjectProperty("c", runtime.NewInt(2)),
		)

		keys, err := objects.Keys(context.Background(), obj)
		keysArray := keys.(*runtime.Array)

		So(err, ShouldEqual, nil)
		actualLength, _ := keysArray.Length(context.Background())
		So(actualLength, ShouldEqual, 3)

		// Check that all expected keys are present (order doesn't matter)
		keyStrings := make([]string, 0, 3)
		keysArray.ForEach(context.Background(), func(ctx context.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			keyStrings = append(keyStrings, val.String())
			return true, nil
		})
		
		for _, expectedKey := range []string{"a", "b", "c"} {
			found := false
			for _, actualKey := range keyStrings {
				if actualKey == expectedKey {
					found = true
					break
				}
			}
			So(found, ShouldBeTrue)
		}
	})

	Convey("Keys(obj, true) should return ['a', 'b', 'c'] in sorted order", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("b", runtime.NewInt(0)),
			runtime.NewObjectProperty("a", runtime.NewInt(1)),
			runtime.NewObjectProperty("c", runtime.NewInt(3)),
		)

		keys, err := objects.Keys(context.Background(), obj, runtime.NewBoolean(true))
		keysArray := keys.(*runtime.Array)

		So(err, ShouldEqual, nil)

		expectedKeys := []string{"a", "b", "c"}
		for idx, key := range expectedKeys {
			actualKey, _ := keysArray.Get(context.Background(), runtime.NewInt(idx))
			So(actualKey.String(), ShouldEqual, key)
		}
	})

	Convey("When there are no keys", t, func() {
		obj := runtime.NewObject()

		keys, err := objects.Keys(context.Background(), obj, runtime.NewBoolean(true))
		keysArray := keys.(*runtime.Array)

		So(err, ShouldEqual, nil)
		actualLength, _ := keysArray.Length(context.Background())
		So(actualLength, ShouldEqual, 0)

		keys, err = objects.Keys(context.Background(), obj, runtime.NewBoolean(false))
		keysArray = keys.(*runtime.Array)

		So(err, ShouldEqual, nil)
		actualLength, _ = keysArray.Length(context.Background())
		So(actualLength, ShouldEqual, 0)
	})

	Convey("When not enough arguments", t, func() {
		_, err := objects.Keys(context.Background())

		So(err, ShouldBeError)
	})

	Convey("When first argument isn't object", t, func() {
		notObj := runtime.NewInt(0)

		_, err := objects.Keys(context.Background(), notObj)

		So(err, ShouldBeError)
	})

	Convey("When second argument isn't boolean", t, func() {
		obj := runtime.NewObject()

		_, err := objects.Keys(context.Background(), obj, obj)

		So(err, ShouldBeError)
	})

	Convey("When object has special character keys", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("key with spaces", runtime.NewInt(1)),
			runtime.NewObjectProperty("key_with_underscores", runtime.NewInt(2)),
			runtime.NewObjectProperty("key-with-dashes", runtime.NewInt(3)),
			runtime.NewObjectProperty("key.with.dots", runtime.NewInt(4)),
		)

		keys, err := objects.Keys(context.Background(), obj, runtime.NewBoolean(true))
		keysArray := keys.(*runtime.Array)

		So(err, ShouldEqual, nil)
		actualLength, _ := keysArray.Length(context.Background())
		So(actualLength, ShouldEqual, 4)

		// Check sorted order
		expectedKeys := []string{"key with spaces", "key-with-dashes", "key.with.dots", "key_with_underscores"}
		for idx, expectedKey := range expectedKeys {
			actualKey, _ := keysArray.Get(context.Background(), runtime.NewInt(idx))
			So(actualKey.String(), ShouldEqual, expectedKey)
		}
	})
}


