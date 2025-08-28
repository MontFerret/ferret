package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeRecursive(t *testing.T) {
	Convey("Wrong arguments", t, func() {
		Convey("It should error when 0 arguments", func() {
			actual, err := objects.MergeRecursive(context.Background())

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, runtime.None), ShouldEqual, 0)
		})

		Convey("It should error when there is not object arguments", func() {
			actual, err := objects.MergeRecursive(context.Background(), runtime.NewInt(0))

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, runtime.None), ShouldEqual, 0)

			actual, err = objects.MergeRecursive(context.Background(),
				runtime.NewInt(0), runtime.NewObject(),
			)

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, runtime.None), ShouldEqual, 0)
		})
	})

	Convey("Merge single object", t, func() {
		obj := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(0)),
		)
		expected := runtime.NewObjectWith(
			runtime.NewObjectProperty("a", runtime.NewInt(0)),
		)

		actual, err := objects.MergeRecursive(context.Background(), obj)

		So(err, ShouldBeNil)
		So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
	})

	Convey("Merge two objects", t, func() {
		Convey("When there are no common keys", func() {
			obj1 := runtime.NewObjectWith(
				runtime.NewObjectProperty("a", runtime.NewInt(0)),
				runtime.NewObjectProperty("b", runtime.NewInt(1)),
			)
			obj2 := runtime.NewObjectWith(
				runtime.NewObjectProperty("c", runtime.NewInt(2)),
				runtime.NewObjectProperty("d", runtime.NewInt(3)),
			)
			expected := runtime.NewObjectWith(
				runtime.NewObjectProperty("a", runtime.NewInt(0)),
				runtime.NewObjectProperty("b", runtime.NewInt(1)),
				runtime.NewObjectProperty("c", runtime.NewInt(2)),
				runtime.NewObjectProperty("d", runtime.NewInt(3)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("When objects with the same key", func() {
			obj1 := runtime.NewObjectWith(
				runtime.NewObjectProperty("a", runtime.NewInt(0)),
				runtime.NewObjectProperty("b", runtime.NewInt(1)),
			)
			obj2 := runtime.NewObjectWith(
				runtime.NewObjectProperty("a", runtime.NewInt(2)), // Same key, different value
				runtime.NewObjectProperty("c", runtime.NewInt(3)),
			)
			expected := runtime.NewObjectWith(
				runtime.NewObjectProperty("a", runtime.NewInt(2)), // Second object's value should win
				runtime.NewObjectProperty("b", runtime.NewInt(1)),
				runtime.NewObjectProperty("c", runtime.NewInt(3)),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)

			So(err, ShouldBeNil)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("Merge two objects with the same keys and nested objects", func() {
			nestedObj1 := runtime.NewObjectWith(
				runtime.NewObjectProperty("x", runtime.NewInt(1)),
			)
			nestedObj2 := runtime.NewObjectWith(
				runtime.NewObjectProperty("y", runtime.NewInt(2)),
			)
			
			obj1 := runtime.NewObjectWith(
				runtime.NewObjectProperty("nested", nestedObj1),
				runtime.NewObjectProperty("simple", runtime.NewString("value1")),
			)
			obj2 := runtime.NewObjectWith(
				runtime.NewObjectProperty("nested", nestedObj2), // Should merge recursively
				runtime.NewObjectProperty("simple", runtime.NewString("value2")), // Should overwrite
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)
			actualObj := actual.(*runtime.Object)

			So(err, ShouldBeNil)

			// Check simple property - should be overwritten by later value
			simpleVal, _ := actualObj.Get(context.Background(), runtime.NewString("simple"))
			So(runtime.CompareValues(simpleVal, runtime.NewString("value2")), ShouldEqual, 0)

			// Check nested object was merged recursively
			nestedVal, _ := actualObj.Get(context.Background(), runtime.NewString("nested"))
			nestedObj := nestedVal.(*runtime.Object)
			
			// Should have both x and y
			xVal, _ := nestedObj.Get(context.Background(), runtime.NewString("x"))
			So(runtime.CompareValues(xVal, runtime.NewInt(1)), ShouldEqual, 0)
			
			yVal, _ := nestedObj.Get(context.Background(), runtime.NewString("y"))
			So(runtime.CompareValues(yVal, runtime.NewInt(2)), ShouldEqual, 0)
		})

		Convey("When there are nested arrays", func() {
			arr1 := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
			arr2 := runtime.NewArrayWith(runtime.NewInt(3), runtime.NewInt(4))
			
			obj1 := runtime.NewObjectWith(
				runtime.NewObjectProperty("arr", arr1),
			)
			obj2 := runtime.NewObjectWith(
				runtime.NewObjectProperty("arr", arr2),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)
			actualObj := actual.(*runtime.Object)

			So(err, ShouldBeNil)

			// Arrays should not be merged recursively, second one should win
			arrVal, _ := actualObj.Get(context.Background(), runtime.NewString("arr"))
			arrResult := arrVal.(*runtime.Array)
			arrLength, _ := arrResult.Length(context.Background())
			So(arrLength, ShouldEqual, 2)
			
			val0, _ := arrResult.Get(context.Background(), runtime.NewInt(0))
			val1, _ := arrResult.Get(context.Background(), runtime.NewInt(1))
			So(runtime.CompareValues(val0, runtime.NewInt(3)), ShouldEqual, 0)
			So(runtime.CompareValues(val1, runtime.NewInt(4)), ShouldEqual, 0)
		})

		Convey("When there are nested objects (example from ArangoDB doc)", func() {
			// { "user-1": { "name": "Jane", "livesIn": { "city": "LA" } } }
			obj1 := runtime.NewObjectWith(
				runtime.NewObjectProperty(
					"user-1", runtime.NewObjectWith(
						runtime.NewObjectProperty(
							"name", runtime.NewString("Jane"),
						),
						runtime.NewObjectProperty(
							"livesIn", runtime.NewObjectWith(
								runtime.NewObjectProperty(
									"city", runtime.NewString("LA"),
								),
							),
						),
					),
				),
			)
			// { "user-1": { "age": 42, "livesIn": { "state": "CA" } } }
			obj2 := runtime.NewObjectWith(
				runtime.NewObjectProperty(
					"user-1", runtime.NewObjectWith(
						runtime.NewObjectProperty(
							"age", runtime.NewInt(42),
						),
						runtime.NewObjectProperty(
							"livesIn", runtime.NewObjectWith(
								runtime.NewObjectProperty(
									"state", runtime.NewString("CA"),
								),
							),
						),
					),
				),
			)

			actual, err := objects.MergeRecursive(context.Background(), obj1, obj2)
			actualObj := actual.(*runtime.Object)

			So(err, ShouldBeNil)

			// Get the merged user object
			userVal, _ := actualObj.Get(context.Background(), runtime.NewString("user-1"))
			userObj := userVal.(*runtime.Object)

			// Should have name from obj1
			nameVal, _ := userObj.Get(context.Background(), runtime.NewString("name"))
			So(runtime.CompareValues(nameVal, runtime.NewString("Jane")), ShouldEqual, 0)

			// Should have age from obj2
			ageVal, _ := userObj.Get(context.Background(), runtime.NewString("age"))
			So(runtime.CompareValues(ageVal, runtime.NewInt(42)), ShouldEqual, 0)

			// Should have nested livesIn object with both city and state
			livesInVal, _ := userObj.Get(context.Background(), runtime.NewString("livesIn"))
			livesInObj := livesInVal.(*runtime.Object)

			cityVal, _ := livesInObj.Get(context.Background(), runtime.NewString("city"))
			So(runtime.CompareValues(cityVal, runtime.NewString("LA")), ShouldEqual, 0)

			stateVal, _ := livesInObj.Get(context.Background(), runtime.NewString("state"))
			So(runtime.CompareValues(stateVal, runtime.NewString("CA")), ShouldEqual, 0)
		})
	})

	Convey("Merged object should be independent of source objects", t, func() {
		Convey("When array", func() {
			arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
			obj := runtime.NewObjectWith(runtime.NewObjectProperty("arr", arr))

			actual, err := objects.MergeRecursive(context.Background(), obj)

			So(err, ShouldBeNil)

			// Modify original array
			arr.Add(context.Background(), runtime.NewInt(3))

			// Get result array and check it's unchanged
			actualObj := actual.(*runtime.Object)
			resultArrVal, _ := actualObj.Get(context.Background(), runtime.NewString("arr"))
			resultArr := resultArrVal.(*runtime.Array)
			resultLength, _ := resultArr.Length(context.Background())
			So(resultLength, ShouldEqual, 2)
		})

		Convey("When object", func() {
			nested := runtime.NewObjectWith(
				runtime.NewObjectProperty("nested", runtime.NewInt(0)),
			)
			obj := runtime.NewObjectWith(runtime.NewObjectProperty("obj", nested))

			actual, err := objects.MergeRecursive(context.Background(), obj)

			So(err, ShouldBeNil)

			// Modify original nested object
			nested.Set(context.Background(), runtime.NewString("str"), runtime.NewInt(0))

			// Get result object and check it's unchanged
			actualObj := actual.(*runtime.Object)
			resultObjVal, _ := actualObj.Get(context.Background(), runtime.NewString("obj"))
			resultObj := resultObjVal.(*runtime.Object)
			resultLength, _ := resultObj.Length(context.Background())
			So(resultLength, ShouldEqual, 1)

			// Should not have the new property
			hasNewKey, _ := resultObj.ContainsKey(context.Background(), runtime.NewString("str"))
			So(hasNewKey, ShouldEqual, runtime.False)
		})
	})
}