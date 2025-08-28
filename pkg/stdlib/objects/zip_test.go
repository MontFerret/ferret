package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"

	. "github.com/smartystreets/goconvey/convey"
)

func TestZip(t *testing.T) {
	Convey("Invalid arguments", t, func() {
		Convey("When there are no arguments", func() {
			actual, err := objects.Zip(context.Background())
			expected := runtime.None

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("When single argument", func() {
			actual, err := objects.Zip(context.Background(), runtime.NewArray(0))
			expected := runtime.None

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)

			actual, err = objects.Zip(context.Background(), runtime.NewInt(0))

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("When too many arguments", func() {
			actual, err := objects.Zip(context.Background(),
				runtime.NewArray(0), runtime.NewArray(0), runtime.NewArray(0))
			expected := runtime.None

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("When there is not array argument", func() {
			actual, err := objects.Zip(context.Background(), runtime.NewArray(0), runtime.NewInt(0))
			expected := runtime.None

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)

			actual, err = objects.Zip(context.Background(), runtime.NewInt(0), runtime.NewArray(0))

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("When there is not string element into keys array", func() {
			keys := runtime.NewArrayWith(runtime.NewInt(0))
			vals := runtime.NewArrayWith(runtime.NewString("v1"))
			expected := runtime.None

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("When keys and values have different lengths", func() {
			keys := runtime.NewArrayWith(runtime.NewString("k1"))
			vals := runtime.NewArrayWith(runtime.NewString("v1"), runtime.NewString("v2"))
			expected := runtime.None

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeError)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})
	})

	Convey("Valid arguments", t, func() {
		Convey("Zip empty arrays", func() {
			keys := runtime.NewArray(0)
			vals := runtime.NewArray(0)
			expected := runtime.NewObject()

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeNil)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("Zip single key-value pair", func() {
			keys := runtime.NewArrayWith(runtime.NewString("k1"))
			vals := runtime.NewArrayWith(runtime.NewInt(1))
			expected := runtime.NewObjectWith(
				runtime.NewObjectProperty("k1", runtime.NewInt(1)),
			)

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeNil)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("Zip multiple key-value pairs", func() {
			keys := runtime.NewArrayWith(
				runtime.NewString("k1"),
				runtime.NewString("k2"),
			)
			vals := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
			)
			expected := runtime.NewObjectWith(
				runtime.NewObjectProperty("k1", runtime.NewInt(1)),
				runtime.NewObjectProperty("k2", runtime.NewInt(2)),
			)

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeNil)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("Zip with duplicate keys - first value wins", func() {
			keys := runtime.NewArrayWith(
				runtime.NewString("a"),
				runtime.NewString("b"),
				runtime.NewString("a"),
			)
			vals := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			expected := runtime.NewObjectWith(
				runtime.NewObjectProperty("a", runtime.NewInt(1)),
				runtime.NewObjectProperty("b", runtime.NewInt(2)),
			)

			actual, err := objects.Zip(context.Background(), keys, vals)

			So(err, ShouldBeNil)
			So(runtime.CompareValues(actual, expected), ShouldEqual, 0)
		})

		Convey("Zip with complex values", func() {
			keys := runtime.NewArrayWith(
				runtime.NewString("arr"),
				runtime.NewString("obj"),
			)
			vals := runtime.NewArrayWith(
				runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
				runtime.NewObjectWith(runtime.NewObjectProperty("nested", runtime.NewString("value"))),
			)

			actual, err := objects.Zip(context.Background(), keys, vals)
			actualObj := actual.(*runtime.Object)

			So(err, ShouldBeNil)

			// Check array value
			arrVal, _ := actualObj.Get(context.Background(), runtime.NewString("arr"))
			arr := arrVal.(*runtime.Array)
			arrLength, _ := arr.Length(context.Background())
			So(arrLength, ShouldEqual, 2)

			// Check object value  
			objVal, _ := actualObj.Get(context.Background(), runtime.NewString("obj"))
			obj := objVal.(*runtime.Object)
			nestedVal, _ := obj.Get(context.Background(), runtime.NewString("nested"))
			So(runtime.CompareValues(nestedVal, runtime.NewString("value")), ShouldEqual, 0)
		})
	})

	Convey("Result is independent of source values", t, func() {
		arr := runtime.NewArrayWith(runtime.NewInt(0))
		keys := runtime.NewArrayWith(runtime.NewString("arr"))
		vals := runtime.NewArrayWith(arr)

		actual, err := objects.Zip(context.Background(), keys, vals)
		actualObj := actual.(*runtime.Object)

		So(err, ShouldBeNil)

		// Modify original array
		arr.Add(context.Background(), runtime.NewInt(1))

		// Check that the result wasn't affected
		resultArr, _ := actualObj.Get(context.Background(), runtime.NewString("arr"))
		resultArrVal := resultArr.(*runtime.Array)
		resultLength, _ := resultArrVal.Length(context.Background())
		So(resultLength, ShouldEqual, 1)

		val, _ := resultArrVal.Get(context.Background(), runtime.NewInt(0))
		So(runtime.CompareValues(val, runtime.NewInt(0)), ShouldEqual, 0)
	})
}