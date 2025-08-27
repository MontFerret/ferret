package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

// TestArgumentValidation tests argument validation across all functions
func TestArgumentValidation(t *testing.T) {
	ctx := context.Background()

	Convey("Argument validation", t, func() {
		Convey("Functions should reject too few arguments", func() {
			// Test functions that require at least 2 arguments
			arr := runtime.NewArrayWith(runtime.NewInt(1))

			_, err := arrays.Append(ctx)
			So(err, ShouldNotBeNil)

			_, err = arrays.Union(ctx, arr)
			So(err, ShouldNotBeNil)

			_, err = arrays.Intersection(ctx, arr)
			So(err, ShouldNotBeNil)

			_, err = arrays.Minus(ctx, arr)
			So(err, ShouldNotBeNil)

			_, err = arrays.RemoveValue(ctx, arr)
			So(err, ShouldNotBeNil)

			_, err = arrays.Position(ctx, arr)
			So(err, ShouldNotBeNil)

			_, err = arrays.Slice(ctx, arr)
			So(err, ShouldNotBeNil)

			_, err = arrays.Nth(ctx, arr)
			So(err, ShouldNotBeNil)
		})

		Convey("Functions should reject invalid argument types", func() {
			// Test passing non-array as first argument
			nonArray := runtime.NewString("not an array")

			_, err := arrays.Sorted(ctx, nonArray)
			So(err, ShouldNotBeNil)

			_, err = arrays.Pop(ctx, nonArray)
			So(err, ShouldNotBeNil)

			_, err = arrays.Shift(ctx, nonArray)
			So(err, ShouldNotBeNil)

			_, err = arrays.First(ctx, nonArray)
			So(err, ShouldNotBeNil)

			_, err = arrays.Last(ctx, nonArray)
			So(err, ShouldNotBeNil)

			// Test passing non-integer as index argument
			arr := runtime.NewArrayWith(runtime.NewInt(1))
			nonInt := runtime.NewString("not an int")

			_, err = arrays.Nth(ctx, arr, nonInt)
			So(err, ShouldNotBeNil)

			_, err = arrays.Slice(ctx, arr, nonInt)
			So(err, ShouldNotBeNil)

			_, err = arrays.RemoveNth(ctx, arr, nonInt)
			So(err, ShouldNotBeNil)

			// Test passing non-boolean as boolean argument
			nonBool := runtime.NewString("not a bool")

			_, err = arrays.Position(ctx, arr, runtime.NewInt(1), nonBool)
			So(err, ShouldNotBeNil)

			_, err = arrays.Append(ctx, arr, runtime.NewInt(1), nonBool)
			So(err, ShouldNotBeNil)
		})
	})
}

// TestDataTypeHandling tests how functions handle different data types
func TestDataTypeHandling(t *testing.T) {
	ctx := context.Background()

	Convey("Data type handling", t, func() {
		Convey("Functions should handle mixed data types", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(42),
				runtime.NewString("hello"),
				runtime.NewBoolean(true),
				runtime.None,
				runtime.NewFloat(3.14),
			)

			// Test that functions work with mixed types
			out, err := arrays.First(ctx, arr)
			So(err, ShouldBeNil)
			So(out.(runtime.Comparable).Compare(runtime.NewInt(42)), ShouldEqual, 0)

			out, err = arrays.Last(ctx, arr)
			So(err, ShouldBeNil)
			So(out.(runtime.Comparable).Compare(runtime.NewFloat(3.14)), ShouldEqual, 0)

			out, err = arrays.Nth(ctx, arr, runtime.NewInt(2))
			So(err, ShouldBeNil)
			So(out.(runtime.Comparable).Compare(runtime.NewBoolean(true)), ShouldEqual, 0)

			// Test Position with different types
			pos, err := arrays.Position(ctx, arr, runtime.None)
			So(err, ShouldBeNil)
			So(pos.String(), ShouldEqual, "true")

			pos, err = arrays.Position(ctx, arr, runtime.NewString("hello"), runtime.NewBoolean(true))
			So(err, ShouldBeNil)
			So(pos.String(), ShouldEqual, "1")
		})

		Convey("Sorted should handle different comparable types separately", func() {
			// Numbers should sort numerically
			numbers := runtime.NewArrayWith(
				runtime.NewInt(3),
				runtime.NewInt(1),
				runtime.NewInt(2),
			)
			out, err := arrays.Sorted(ctx, numbers)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[1,2,3]")

			// Strings should sort alphabetically
			strings := runtime.NewArrayWith(
				runtime.NewString("c"),
				runtime.NewString("a"),
				runtime.NewString("b"),
			)
			out, err = arrays.Sorted(ctx, strings)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, `["a","b","c"]`)
		})
	})
}

// TestPerformanceScenarios tests functions with larger datasets
func TestPerformanceScenarios(t *testing.T) {
	ctx := context.Background()

	Convey("Performance scenarios", t, func() {
		Convey("Should handle large arrays efficiently", func() {
			// Create a large array with 1000 elements
			arr := runtime.NewArray(1000)
			for i := 0; i < 1000; i++ {
				arr.Add(ctx, runtime.NewInt(i))
			}

			// Test operations that should work efficiently
			out, err := arrays.First(ctx, arr)
			So(err, ShouldBeNil)
			So(out.(runtime.Comparable).Compare(runtime.NewInt(0)), ShouldEqual, 0)

			out, err = arrays.Last(ctx, arr)
			So(err, ShouldBeNil)
			So(out.(runtime.Comparable).Compare(runtime.NewInt(999)), ShouldEqual, 0)

			out, err = arrays.Nth(ctx, arr, runtime.NewInt(500))
			So(err, ShouldBeNil)
			So(out.(runtime.Comparable).Compare(runtime.NewInt(500)), ShouldEqual, 0)

			// Test slice operations
			out, err = arrays.Slice(ctx, arr, runtime.NewInt(100), runtime.NewInt(10))
			So(err, ShouldBeNil)
			length, lengthErr := out.(runtime.List).Length(ctx)
			So(lengthErr, ShouldBeNil)
			So(length, ShouldEqual, 10)
		})

		Convey("Should handle empty arrays gracefully", func() {
			emptyArr := runtime.NewArrayWith()

			out, err := arrays.First(ctx, emptyArr)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.None)

			out, err = arrays.Last(ctx, emptyArr)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.None)

			out, err = arrays.Pop(ctx, emptyArr)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[]")

			out, err = arrays.Shift(ctx, emptyArr)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[]")

			out, err = arrays.Sorted(ctx, emptyArr)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[]")
		})
	})
}