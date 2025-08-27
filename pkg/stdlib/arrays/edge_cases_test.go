package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

// TestEdgeCases tests various edge cases across all array functions
func TestEdgeCases(t *testing.T) {
	ctx := context.Background()

	Convey("Flatten edge cases", t, func() {
		Convey("Should handle empty arrays", func() {
			arr := runtime.NewArrayWith()
			out, err := arrays.Flatten(ctx, arr)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[]")
		})

		Convey("Should handle depth 0 correctly", func() {
			arr := runtime.NewArrayWith(
				runtime.NewArrayWith(runtime.NewInt(1)),
				runtime.NewArrayWith(runtime.NewInt(2)),
			)
			out, err := arrays.Flatten(ctx, arr, runtime.NewInt(0))
			So(err, ShouldBeNil)
			// With depth 0, no flattening should occur
			So(out.String(), ShouldEqual, "[[1],[2]]")
		})

		Convey("Should handle negative depth", func() {
			arr := runtime.NewArrayWith(
				runtime.NewArrayWith(runtime.NewInt(1)),
				runtime.NewArrayWith(runtime.NewInt(2)),
			)
			out, err := arrays.Flatten(ctx, arr, runtime.NewInt(-1))
			So(err, ShouldBeNil)
			// Negative depth should not flatten
			So(out.String(), ShouldEqual, "[[1],[2]]")
		})

		Convey("Should handle very deep nesting", func() {
			// Create [[[[[5]]]]]
			deep := runtime.Value(runtime.NewInt(5))
			for i := 0; i < 5; i++ {
				deep = runtime.NewArrayWith(deep)
			}
			arr := runtime.NewArrayWith(deep)

			out, err := arrays.Flatten(ctx, arr, runtime.NewInt(6))
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[5]")
		})

		Convey("Should handle mixed types correctly", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewArrayWith(
					runtime.NewString("test"),
					runtime.NewBoolean(true),
				),
				runtime.NewString("end"),
			)
			out, err := arrays.Flatten(ctx, arr)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, `[1,"test",true,"end"]`)
		})
	})

	Convey("RemoveValue edge cases", t, func() {
		Convey("Should handle limit of 0", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(0))
			So(err, ShouldBeNil)
			// With limit 0, nothing should be removed
			So(out.String(), ShouldEqual, "[1,2,2,3]")
		})

		Convey("Should handle negative limit", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(-5))
			So(err, ShouldBeNil)
			// Negative limit should remove all occurrences
			So(out.String(), ShouldEqual, "[1,3]")
		})

		Convey("Should handle removing non-existent value", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(5))
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[1,2,3]")
		})
	})

	Convey("Slice edge cases", t, func() {
		Convey("Should handle negative start index", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			_, err := arrays.Slice(ctx, arr, runtime.NewInt(-1))
			// Implementation should handle this gracefully - may return error or empty array
			// Just check that it doesn't panic
			So(err, ShouldBeNil)
		})

		Convey("Should handle negative length", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			_, err := arrays.Slice(ctx, arr, runtime.NewInt(1), runtime.NewInt(-1))
			// Implementation should handle negative length
			// Just check that it doesn't panic
			So(err, ShouldBeNil)
		})

		Convey("Should handle zero length", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			out, err := arrays.Slice(ctx, arr, runtime.NewInt(1), runtime.NewInt(0))
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[]")
		})
	})

	Convey("Nth edge cases", t, func() {
		Convey("Should handle very large index", func() {
			arr := runtime.NewArrayWith(runtime.NewInt(1))
			out, err := arrays.Nth(ctx, arr, runtime.NewInt(1000))
			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.None)
		})

		Convey("Should handle negative index correctly", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			out, err := arrays.Nth(ctx, arr, runtime.NewInt(-1))
			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.None)
		})
	})

	Convey("Position edge cases", t, func() {
		Convey("Should handle searching for None", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.None,
				runtime.NewInt(3),
			)
			out, err := arrays.Position(ctx, arr, runtime.None)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "true")
		})

		Convey("Should handle empty array", func() {
			arr := runtime.NewArrayWith()
			out, err := arrays.Position(ctx, arr, runtime.NewInt(1))
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "false")
		})

		Convey("Should return correct index with position flag", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)
			out, err := arrays.Position(ctx, arr, runtime.NewInt(2), runtime.NewBoolean(true))
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "1")
		})
	})

	Convey("Intersection edge cases", t, func() {
		Convey("Should handle empty arrays", func() {
			arr1 := runtime.NewArrayWith()
			arr2 := runtime.NewArrayWith(runtime.NewInt(1))
			out, err := arrays.Intersection(ctx, arr1, arr2)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[]")
		})

		Convey("Should handle identical arrays", func() {
			arr1 := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
			arr2 := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
			out, err := arrays.Intersection(ctx, arr1, arr2)
			So(err, ShouldBeNil)
			// Should contain both elements
			length, lengthErr := out.(runtime.List).Length(ctx)
			So(lengthErr, ShouldBeNil)
			So(length, ShouldEqual, 2)
		})
	})

	Convey("Union edge cases", t, func() {
		Convey("Should handle empty arrays", func() {
			arr1 := runtime.NewArrayWith()
			arr2 := runtime.NewArrayWith()
			out, err := arrays.Union(ctx, arr1, arr2)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[]")
		})
	})
}