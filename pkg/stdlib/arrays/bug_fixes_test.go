package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
)

// TestBugFixes tests for specific bugs that were found and fixed
func TestBugFixes(t *testing.T) {
	ctx := context.Background()

	Convey("Bug fixes", t, func() {
		Convey("RemoveValue should handle edge cases correctly", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)

			Convey("Limit of 0 should remove nothing", func() {
				out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(0))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[1,2,2,3]")
			})

			Convey("Negative limit should remove all occurrences", func() {
				out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(-1))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[1,3]")

				out, err = arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(-10))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[1,3]")
			})

			Convey("Positive limit should remove up to limit", func() {
				out, err := arrays.RemoveValue(ctx, arr, runtime.NewInt(2), runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[1,2,3]")
			})
		})

		Convey("Slice should handle edge cases without panicking", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.NewInt(2),
				runtime.NewInt(3),
			)

			Convey("Negative start index should return empty array", func() {
				out, err := arrays.Slice(ctx, arr, runtime.NewInt(-1))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[]")

				out, err = arrays.Slice(ctx, arr, runtime.NewInt(-10))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[]")
			})

			Convey("Negative length should return empty array", func() {
				out, err := arrays.Slice(ctx, arr, runtime.NewInt(1), runtime.NewInt(-1))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[]")

				out, err = arrays.Slice(ctx, arr, runtime.NewInt(0), runtime.NewInt(-5))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[]")
			})

			Convey("Start index beyond array should return empty array", func() {
				out, err := arrays.Slice(ctx, arr, runtime.NewInt(10))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[]")
			})

			Convey("Length beyond array should return up to end", func() {
				out, err := arrays.Slice(ctx, arr, runtime.NewInt(1), runtime.NewInt(100))
				So(err, ShouldBeNil)
				So(out.String(), ShouldEqual, "[2,3]")
			})
		})

		Convey("Nth should handle edge cases without panicking", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(10),
				runtime.NewInt(20),
				runtime.NewInt(30),
			)

			Convey("Negative index should return None", func() {
				out, err := arrays.Nth(ctx, arr, runtime.NewInt(-1))
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.None)

				out, err = arrays.Nth(ctx, arr, runtime.NewInt(-100))
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.None)
			})

			Convey("Index beyond array should return None", func() {
				out, err := arrays.Nth(ctx, arr, runtime.NewInt(3))
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.None)

				out, err = arrays.Nth(ctx, arr, runtime.NewInt(1000))
				So(err, ShouldBeNil)
				So(out, ShouldEqual, runtime.None)
			})

			Convey("Valid indices should work correctly", func() {
				out, err := arrays.Nth(ctx, arr, runtime.NewInt(0))
				So(err, ShouldBeNil)
				So(out.(runtime.Comparable).Compare(runtime.NewInt(10)), ShouldEqual, 0)

				out, err = arrays.Nth(ctx, arr, runtime.NewInt(2))
				So(err, ShouldBeNil)
				So(out.(runtime.Comparable).Compare(runtime.NewInt(30)), ShouldEqual, 0)
			})
		})

		Convey("Last should properly validate arguments", func() {
			nonArray := runtime.NewString("not an array")

			Convey("Should return error for non-array input", func() {
				_, err := arrays.Last(ctx, nonArray)
				So(err, ShouldNotBeNil)
			})

			Convey("Should work correctly with valid array", func() {
				arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
				out, err := arrays.Last(ctx, arr)
				So(err, ShouldBeNil)
				So(out.(runtime.Comparable).Compare(runtime.NewInt(2)), ShouldEqual, 0)
			})
		})
	})
}

// TestSpecialValues tests functions with special runtime values
func TestSpecialValues(t *testing.T) {
	ctx := context.Background()

	Convey("Special values handling", t, func() {
		Convey("Functions should handle None values correctly", func() {
			arr := runtime.NewArrayWith(
				runtime.NewInt(1),
				runtime.None,
				runtime.NewInt(3),
			)

			// Position should find None
			pos, err := arrays.Position(ctx, arr, runtime.None)
			So(err, ShouldBeNil)
			So(pos.String(), ShouldEqual, "true")

			// Position should return correct index for None
			pos, err = arrays.Position(ctx, arr, runtime.None, runtime.NewBoolean(true))
			So(err, ShouldBeNil)
			So(pos.String(), ShouldEqual, "1")

			// RemoveValue should remove None
			out, err := arrays.RemoveValue(ctx, arr, runtime.None)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[1,3]")

			// Nth should return None when accessing None
			noneVal, err := arrays.Nth(ctx, arr, runtime.NewInt(1))
			So(err, ShouldBeNil)
			So(noneVal, ShouldEqual, runtime.None)
		})

		Convey("Functions should handle arrays with only None values", func() {
			noneOnlyArr := runtime.NewArrayWith(runtime.None, runtime.None, runtime.None)

			out, err := arrays.First(ctx, noneOnlyArr)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.None)

			out, err = arrays.Last(ctx, noneOnlyArr)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.None)

			out, err = arrays.Unique(ctx, noneOnlyArr)
			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "[null]")
		})
	})
}