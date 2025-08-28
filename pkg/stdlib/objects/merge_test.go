package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMerge(t *testing.T) {
	Convey("When not enough arguments", t, func() {
		obj, err := objects.Merge(context.Background())

		So(err, ShouldBeError)
		So(runtime.CompareValues(obj, runtime.None), ShouldEqual, 0)
	})

	Convey("When wrong type of arguments", t, func() {
		obj, err := objects.Merge(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeError)
		So(runtime.CompareValues(obj, runtime.None), ShouldEqual, 0)

		obj, err = objects.Merge(context.Background(), runtime.NewObject(), runtime.NewInt(0))

		So(err, ShouldBeError)
		So(runtime.CompareValues(obj, runtime.None), ShouldEqual, 0)
	})

	Convey("When array contains non-objects", t, func() {
		arr := runtime.NewArrayWith(runtime.NewObject(), runtime.NewInt(0))
		obj, err := objects.Merge(context.Background(), arr)

		So(err, ShouldBeError)
		So(runtime.CompareValues(obj, runtime.None), ShouldEqual, 0)
	})

	Convey("Merged object should be independent of source objects", t, func() {
		obj1 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
		)
		obj2 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop3", runtime.NewInt(3)),
		)

		result := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
			runtime.NewObjectProperty("prop3", runtime.NewInt(3)),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)

		So(err, ShouldBeNil)
		So(runtime.CompareValues(merged, result), ShouldEqual, 0)

		// Modify original objects to ensure independence
		obj1.Set(context.Background(), runtime.NewString("newProp"), runtime.NewString("newVal"))
		obj2.Set(context.Background(), runtime.NewString("newProp2"), runtime.NewString("newVal2"))

		// Merged object should remain unchanged
		So(runtime.CompareValues(merged, result), ShouldEqual, 0)
	})
}

func TestMergeObjects(t *testing.T) {
	Convey("Merge single object", t, func() {
		obj1 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
		)
		result := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
		)

		merged, err := objects.Merge(context.Background(), obj1)

		So(err, ShouldBeNil)
		So(runtime.CompareValues(merged, result), ShouldEqual, 0)
	})

	Convey("Merge two objects", t, func() {
		obj1 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
		)
		obj2 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop3", runtime.NewInt(3)),
		)

		result := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
			runtime.NewObjectProperty("prop3", runtime.NewInt(3)),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)

		So(err, ShouldBeNil)
		So(runtime.CompareValues(merged, result), ShouldEqual, 0)
	})

	Convey("When keys are repeated", t, func() {
		obj1 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
		)
		obj2 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(2)), // Same key, different value
		)

		// Later objects should overwrite earlier ones
		result := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(2)), // Should be the value from obj2
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)

		So(err, ShouldBeNil)
		So(runtime.CompareValues(merged, result), ShouldEqual, 0)
	})

	Convey("Merge empty objects", t, func() {
		obj1 := runtime.NewObject()
		obj2 := runtime.NewObject()
		result := runtime.NewObject()

		merged, err := objects.Merge(context.Background(), obj1, obj2)

		So(err, ShouldBeNil)
		So(runtime.CompareValues(merged, result), ShouldEqual, 0)
	})

	Convey("Merge objects with complex values", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
		nestedObj := runtime.NewObjectWith(runtime.NewObjectProperty("nested", runtime.NewString("value")))
		
		obj1 := runtime.NewObjectWith(
			runtime.NewObjectProperty("array", arr),
		)
		obj2 := runtime.NewObjectWith(
			runtime.NewObjectProperty("object", nestedObj),
		)

		merged, err := objects.Merge(context.Background(), obj1, obj2)
		mergedObj := merged.(*runtime.Object)

		So(err, ShouldBeNil)

		// Check that array is present
		arrVal, _ := mergedObj.Get(context.Background(), runtime.NewString("array"))
		arrResult := arrVal.(*runtime.Array)
		arrLength, _ := arrResult.Length(context.Background())
		So(arrLength, ShouldEqual, 2)

		// Check that object is present
		objVal, _ := mergedObj.Get(context.Background(), runtime.NewString("object"))
		objResult := objVal.(*runtime.Object)
		nestedVal, _ := objResult.Get(context.Background(), runtime.NewString("nested"))
		So(runtime.CompareValues(nestedVal, runtime.NewString("value")), ShouldEqual, 0)

		// Verify independence - modify original array
		arr.Add(context.Background(), runtime.NewInt(3))
		
		// Merged array should not be affected
		mergedArrLength, _ := arrResult.Length(context.Background())
		So(mergedArrLength, ShouldEqual, 2)
	})
}

func TestMergeArray(t *testing.T) {
	Convey("Merge array", t, func() {
		obj1 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
		)
		obj2 := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop3", runtime.NewInt(3)),
		)

		objArr := runtime.NewArrayWith(obj1, obj2)
		result := runtime.NewObjectWith(
			runtime.NewObjectProperty("prop1", runtime.NewInt(1)),
			runtime.NewObjectProperty("prop2", runtime.NewString("str")),
			runtime.NewObjectProperty("prop3", runtime.NewInt(3)),
		)

		merged, err := objects.Merge(context.Background(), objArr)

		So(err, ShouldBeNil)
		So(runtime.CompareValues(merged, result), ShouldEqual, 0)
	})

	Convey("Merge empty array", t, func() {
		objArr := runtime.NewArray(0)
		result := runtime.NewObject()

		merged, err := objects.Merge(context.Background(), objArr)

		So(err, ShouldBeNil)
		So(runtime.CompareValues(merged, result), ShouldEqual, 0)
	})

	Convey("When there is not object element inside the array", t, func() {
		objArr := runtime.NewArrayWith(
			runtime.NewObject(),
			runtime.NewInt(0),
		)

		obj, err := objects.Merge(context.Background(), objArr)

		So(err, ShouldBeError)
		So(runtime.CompareValues(obj, runtime.None), ShouldEqual, 0)
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