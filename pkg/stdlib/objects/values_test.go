package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValues(t *testing.T) {
	Convey("Invalid arguments", t, func() {
		Convey("When there is no arguments", func() {
			actual, err := objects.Values(context.Background())

			So(err, ShouldBeError)
			So(actual, ShouldEqual, runtime.None)
		})

		Convey("When 2 arguments", func() {
			obj := runtime.NewObjectWith(
				map[string]runtime.Value{
					"k1": runtime.NewInt(0),
					"k2": runtime.NewInt(1),
				},
			)

			actual, err := objects.Values(context.Background(), obj, obj)

			So(err, ShouldBeError)
			So(actual, ShouldEqual, runtime.None)

			actual, err = objects.Values(context.Background(), obj, runtime.NewInt(0))

			So(err, ShouldBeError)
			So(actual, ShouldEqual, runtime.None)
		})

		Convey("When there is not object argument", func() {
			actual, err := objects.Values(context.Background(), runtime.NewInt(0))

			So(err, ShouldBeError)
			So(actual, ShouldEqual, runtime.None)
		})
	})

	Convey("When simple type attributes (same type)", t, func() {
		obj := runtime.NewObjectWith(
			map[string]runtime.Value{
				"k1": runtime.NewInt(0),
				"k2": runtime.NewInt(1),
			},
		)

		actual, err := objects.Values(context.Background(), obj)
		actualArray := actual.(*runtime.Array)

		So(err, ShouldBeNil)
		actualLength, _ := actualArray.Length(context.Background())
		So(actualLength, ShouldEqual, 2)

		// Check that all values are present (order may vary)
		values := make([]runtime.Value, 0, 2)
		actualArray.ForEach(context.Background(), func(ctx context.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			values = append(values, val)
			return true, nil
		})

		// Check that we have both values
		hasZero := false
		hasOne := false
		for _, val := range values {
			if runtime.CompareValues(val, runtime.NewInt(0)) == 0 {
				hasZero = true
			}
			if runtime.CompareValues(val, runtime.NewInt(1)) == 0 {
				hasOne = true
			}
		}
		So(hasZero, ShouldBeTrue)
		So(hasOne, ShouldBeTrue)
	})

	Convey("When simple type attributes (different types)", t, func() {
		obj := runtime.NewObjectWith(
			map[string]runtime.Value{
				"k1": runtime.NewInt(0),
				"k2": runtime.NewString("v2"),
			},
		)

		actual, err := objects.Values(context.Background(), obj)
		actualArray := actual.(*runtime.Array)

		So(err, ShouldBeNil)
		actualLength, _ := actualArray.Length(context.Background())
		So(actualLength, ShouldEqual, 2)

		// Check that all values are present
		values := make([]runtime.Value, 0, 2)
		actualArray.ForEach(context.Background(), func(ctx context.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			values = append(values, val)
			return true, nil
		})

		hasInt := false
		hasString := false
		for _, val := range values {
			if runtime.CompareValues(val, runtime.NewInt(0)) == 0 {
				hasInt = true
			}
			if runtime.CompareValues(val, runtime.NewString("v2")) == 0 {
				hasString = true
			}
		}
		So(hasInt, ShouldBeTrue)
		So(hasString, ShouldBeTrue)
	})

	Convey("When complex type attributes (array)", t, func() {
		arr1 := runtime.NewArrayWith(
			runtime.NewInt(0), runtime.NewInt(1),
		)
		arr2 := runtime.NewArrayWith(
			runtime.NewInt(2), runtime.NewInt(3),
		)
		obj := runtime.NewObjectWith(
			map[string]runtime.Value{
				"k1": arr1,
				"k2": arr2,
			},
		)

		actual, err := objects.Values(context.Background(), obj)
		actualArray := actual.(*runtime.Array)

		So(err, ShouldBeNil)
		actualLength, _ := actualArray.Length(context.Background())
		So(actualLength, ShouldEqual, 2)

		// Check that both arrays are present
		values := make([]runtime.Value, 0, 2)
		actualArray.ForEach(context.Background(), func(ctx context.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			values = append(values, val)
			return true, nil
		})

		hasArr1 := false
		hasArr2 := false
		for _, val := range values {
			if arr, ok := val.(*runtime.Array); ok {
				length, _ := arr.Length(context.Background())
				if length == 2 {
					val0, _ := arr.At(context.Background(), runtime.NewInt(0))
					val1, _ := arr.At(context.Background(), runtime.NewInt(1))
					if runtime.CompareValues(val0, runtime.NewInt(0)) == 0 && runtime.CompareValues(val1, runtime.NewInt(1)) == 0 {
						hasArr1 = true
					}
					if runtime.CompareValues(val0, runtime.NewInt(2)) == 0 && runtime.CompareValues(val1, runtime.NewInt(3)) == 0 {
						hasArr2 = true
					}
				}
			}
		}
		So(hasArr1, ShouldBeTrue)
		So(hasArr2, ShouldBeTrue)
	})

	Convey("When both type attributes", t, func() {
		obj1 := runtime.NewObjectWith(
			map[string]runtime.Value{
				"k1": runtime.NewInt(0),
			},
		)
		arr1 := runtime.NewArrayWith(
			runtime.NewInt(0), runtime.NewInt(1),
		)
		int1 := runtime.NewInt(0)
		obj := runtime.NewObjectWith(
			map[string]runtime.Value{
				"obj": obj1,
				"arr": arr1,
				"int": int1,
			},
		)

		actual, err := objects.Values(context.Background(), obj)
		actualArray := actual.(*runtime.Array)

		So(err, ShouldBeNil)
		actualLength, _ := actualArray.Length(context.Background())
		So(actualLength, ShouldEqual, 3)

		// Check that all three values are present
		hasObject := false
		hasArray := false
		hasInt := false

		actualArray.ForEach(context.Background(), func(ctx context.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			switch val.(type) {
			case *runtime.Object:
				hasObject = true
			case *runtime.Array:
				hasArray = true
			case runtime.Int:
				hasInt = true
			}
			return true, nil
		})

		So(hasObject, ShouldBeTrue)
		So(hasArray, ShouldBeTrue)
		So(hasInt, ShouldBeTrue)
	})

	Convey("Result is independent on the source object (array)", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(0))
		obj := runtime.NewObjectWith(
			map[string]runtime.Value{
				"arr": arr,
			},
		)

		actual, err := objects.Values(context.Background(), obj)
		actualArray := actual.(*runtime.Array)

		So(err, ShouldBeNil)

		// At the returned array value
		returnedArr, _ := actualArray.At(context.Background(), runtime.NewInt(0))
		returnedArrayVal := returnedArr.(*runtime.Array)

		// Modify the original array
		arr.Append(context.Background(), runtime.NewInt(1))

		// Check that the returned array wasn't affected
		returnedLength, _ := returnedArrayVal.Length(context.Background())
		So(returnedLength, ShouldEqual, 1)

		val, _ := returnedArrayVal.At(context.Background(), runtime.NewInt(0))
		So(runtime.CompareValues(val, runtime.NewInt(0)), ShouldEqual, 0)
	})

	Convey("Result is independent on the source object (object)", t, func() {
		nested := runtime.NewObjectWith(
			map[string]runtime.Value{
				"int": runtime.NewInt(0),
			},
		)
		obj := runtime.NewObjectWith(
			map[string]runtime.Value{
				"nested": nested,
			},
		)

		actual, err := objects.Values(context.Background(), obj)
		actualArray := actual.(*runtime.Array)

		So(err, ShouldBeNil)

		// At the returned object value
		returnedObj, _ := actualArray.At(context.Background(), runtime.NewInt(0))
		returnedObjectVal := returnedObj.(*runtime.Object)

		// Modify the original nested object
		nested.Set(context.Background(), runtime.NewString("new"), runtime.NewInt(1))

		// Check that the returned object wasn't affected
		returnedLength, _ := returnedObjectVal.Length(context.Background())
		So(returnedLength, ShouldEqual, 1)

		hasNewKey, _ := returnedObjectVal.ContainsKey(context.Background(), runtime.NewString("new"))
		So(hasNewKey, ShouldEqual, runtime.False)
	})

	Convey("Empty object", t, func() {
		obj := runtime.NewObject()

		actual, err := objects.Values(context.Background(), obj)
		actualArray := actual.(*runtime.Array)

		So(err, ShouldBeNil)
		actualLength, _ := actualArray.Length(context.Background())
		So(actualLength, ShouldEqual, 0)
	})

	Convey("When object has nil values", t, func() {
		obj := runtime.NewObjectWith(
			map[string]runtime.Value{
				"none":         runtime.None,
				"bool":         runtime.NewBoolean(false),
				"empty_string": runtime.NewString(""),
			},
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

	Convey("When object has deeply nested structures", t, func() {
		deepObject := runtime.NewObjectWith(
			map[string]runtime.Value{
				"level1": runtime.NewObjectWith(
					map[string]runtime.Value{
						"level2": runtime.NewObjectWith(
							map[string]runtime.Value{
								"level3": runtime.NewString("deep"),
							},
						),
					}),
			},
		)

		obj := runtime.NewObjectWith(
			map[string]runtime.Value{
				"deep": deepObject,
			},
		)

		actual, err := objects.Values(context.Background(), obj)
		actualArray := actual.(*runtime.Array)

		So(err, ShouldBeNil)
		actualLength, _ := actualArray.Length(context.Background())
		So(actualLength, ShouldEqual, 1)

		// At the deep object and verify it's independent
		returnedDeep, _ := actualArray.At(context.Background(), runtime.NewInt(0))
		returnedDeepObj := returnedDeep.(*runtime.Object)

		// Modify the original deep object
		deepObject.Set(context.Background(), runtime.NewString("modified"), runtime.NewString("value"))

		// Check that the returned object wasn't affected
		hasModified, _ := returnedDeepObj.ContainsKey(context.Background(), runtime.NewString("modified"))
		So(hasModified, ShouldEqual, runtime.False)
	})

	Convey("Stress test", t, func() {
		for i := 0; i < 100; i++ {
			obj1 := runtime.NewObjectWith(
				map[string]runtime.Value{
					"int0": runtime.NewInt(0),
				},
			)
			obj2 := runtime.NewObjectWith(
				map[string]runtime.Value{
					"int1": runtime.NewInt(1),
				},
			)
			obj := runtime.NewObjectWith(
				map[string]runtime.Value{
					"k1": obj1,
					"k2": obj2,
				},
			)

			actual, err := objects.Values(context.Background(), obj)
			actualArray := actual.(*runtime.Array)

			So(err, ShouldBeNil)
			actualLength, _ := actualArray.Length(context.Background())
			So(actualLength, ShouldEqual, 2)

			// Check that both objects are present in the result
			objectCount := 0
			actualArray.ForEach(context.Background(), func(ctx context.Context, val runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
				if _, ok := val.(*runtime.Object); ok {
					objectCount++
				}
				return true, nil
			})
			So(objectCount, ShouldEqual, 2)
		}
	})
}
