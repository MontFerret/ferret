package objects_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHasEdgeCases(t *testing.T) {
	Convey("Edge cases for Has function", t, func() {
		Convey("When key is empty string", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("", runtime.NewString("empty")),
			)

			val, err := objects.Has(context.Background(), obj, runtime.NewString(""))
			valBool := val.(runtime.Boolean)

			So(err, ShouldEqual, nil)
			So(bool(valBool), ShouldEqual, true)
		})

		Convey("When object has many keys", func() {
			properties := make([]*runtime.ObjectProperty, 100)
			for i := 0; i < 100; i++ {
				properties[i] = runtime.NewObjectProperty(
					fmt.Sprintf("key%d", i),
					runtime.NewInt(i),
				)
			}
			obj := runtime.NewObjectWith(properties...)

			// Test existing key
			val, err := objects.Has(context.Background(), obj, runtime.NewString("key50"))
			valBool := val.(runtime.Boolean)

			So(err, ShouldEqual, nil)
			So(bool(valBool), ShouldEqual, true)

			// Test non-existing key
			val, err = objects.Has(context.Background(), obj, runtime.NewString("key999"))
			valBool = val.(runtime.Boolean)

			So(err, ShouldEqual, nil)
			So(bool(valBool), ShouldEqual, false)
		})
	})
}

func TestKeysEdgeCases(t *testing.T) {
	Convey("Edge cases for Keys function", t, func() {
		Convey("When object has special character keys", func() {
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
	})
}

func TestValuesEdgeCases(t *testing.T) {
	Convey("Edge cases for Values function", t, func() {
		Convey("When object has nil values", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("none", runtime.None),
				runtime.NewObjectProperty("bool", runtime.NewBoolean(false)),
				runtime.NewObjectProperty("empty_string", runtime.NewString("")),
			)

			actual, err := objects.Values(context.Background(), obj)
			actualArray := actual.(*runtime.Array)

			So(err, ShouldBeNil)
			actualLength, _ := actualArray.Length(context.Background())
			So(actualLength, ShouldEqual, 3)

			// Verify all values are present
			foundNone := false
			foundBool := false  
			foundEmptyString := false

			actualArray.ForEach(context.Background(), func(ctx context.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
				if runtime.CompareValues(val, runtime.None) == 0 {
					foundNone = true
				}
				if runtime.CompareValues(val, runtime.NewBoolean(false)) == 0 {
					foundBool = true
				}
				if runtime.CompareValues(val, runtime.NewString("")) == 0 {
					foundEmptyString = true
				}
				return true, nil
			})

			So(foundNone, ShouldBeTrue)
			So(foundBool, ShouldBeTrue)
			So(foundEmptyString, ShouldBeTrue)
		})

		Convey("When object has deeply nested structures", func() {
			deepObject := runtime.NewObjectWith(
				runtime.NewObjectProperty("level1", runtime.NewObjectWith(
					runtime.NewObjectProperty("level2", runtime.NewObjectWith(
						runtime.NewObjectProperty("level3", runtime.NewString("deep")),
					)),
				)),
			)

			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("deep", deepObject),
			)

			actual, err := objects.Values(context.Background(), obj)
			actualArray := actual.(*runtime.Array)

			So(err, ShouldBeNil)
			actualLength, _ := actualArray.Length(context.Background())
			So(actualLength, ShouldEqual, 1)

			// Get the deep object and verify it's independent
			returnedDeep, _ := actualArray.Get(context.Background(), runtime.NewInt(0))
			returnedDeepObj := returnedDeep.(*runtime.Object)

			// Modify the original deep object
			deepObject.Set(context.Background(), runtime.NewString("modified"), runtime.NewString("value"))

			// Check that the returned object wasn't affected
			hasModified, _ := returnedDeepObj.ContainsKey(context.Background(), runtime.NewString("modified"))
			So(hasModified, ShouldEqual, runtime.False)
		})
	})
}

func TestZipEdgeCases(t *testing.T) {
	Convey("Edge cases for Zip function", t, func() {
		Convey("When keys contain special characters", func() {
			keys := runtime.NewArrayWith(
				runtime.NewString("key with spaces"),
				runtime.NewString("key_with_underscores"),
				runtime.NewString("key-with-dashes"),
			)
			vals := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)

			actual, err := objects.Zip(context.Background(), keys, vals)
			actualObj := actual.(*runtime.Object)

			So(err, ShouldBeNil)

			// Check all keys are present
			val1, _ := actualObj.Get(context.Background(), runtime.NewString("key with spaces"))
			So(runtime.CompareValues(val1, runtime.NewInt(1)), ShouldEqual, 0)

			val2, _ := actualObj.Get(context.Background(), runtime.NewString("key_with_underscores"))
			So(runtime.CompareValues(val2, runtime.NewInt(2)), ShouldEqual, 0)

			val3, _ := actualObj.Get(context.Background(), runtime.NewString("key-with-dashes"))
			So(runtime.CompareValues(val3, runtime.NewInt(3)), ShouldEqual, 0)
		})

		Convey("When values are all None", func() {
			keys := runtime.NewArrayWith(
				runtime.NewString("k1"),
				runtime.NewString("k2"),
			)
			vals := runtime.NewArrayWith(
				runtime.None,
				runtime.None,
			)

			actual, err := objects.Zip(context.Background(), keys, vals)
			actualObj := actual.(*runtime.Object)

			So(err, ShouldBeNil)

			val1, _ := actualObj.Get(context.Background(), runtime.NewString("k1"))
			So(runtime.CompareValues(val1, runtime.None), ShouldEqual, 0)

			val2, _ := actualObj.Get(context.Background(), runtime.NewString("k2"))
			So(runtime.CompareValues(val2, runtime.None), ShouldEqual, 0)
		})
	})
}

func TestMergeEdgeCases(t *testing.T) {
	Convey("Edge cases for Merge functions", t, func() {
		Convey("Merge with empty objects", func() {
			obj1 := runtime.NewObject()
			obj2 := runtime.NewObjectWith(
				runtime.NewObjectProperty("key", runtime.NewString("value")),
			)

			merged, err := objects.Merge(context.Background(), obj1, obj2)
			mergedObj := merged.(*runtime.Object)

			So(err, ShouldBeNil)

			val, _ := mergedObj.Get(context.Background(), runtime.NewString("key"))
			So(runtime.CompareValues(val, runtime.NewString("value")), ShouldEqual, 0)
		})

		Convey("MergeRecursive with identical objects", func() {
			obj := runtime.NewObjectWith(
				runtime.NewObjectProperty("key", runtime.NewString("value")),
			)

			merged, err := objects.MergeRecursive(context.Background(), obj, obj)

			So(err, ShouldBeNil)
			So(runtime.CompareValues(merged, obj), ShouldEqual, 0)
		})
	})
}