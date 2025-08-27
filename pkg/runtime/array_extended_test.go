package runtime_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func TestArrayExtended(t *testing.T) {
	ctx := context.Background()

	Convey("Array Extended Functionality", t, func() {

		Convey(".Contains", func() {
			Convey("Should return true when element exists", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				contains, err := arr.Contains(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.True)

				contains, err = arr.Contains(ctx, runtime.NewString("test"))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.True)

				contains, err = arr.Contains(ctx, runtime.NewFloat(3.14))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.True)
			})

			Convey("Should return false when element does not exist", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
				)

				contains, err := arr.Contains(ctx, runtime.NewInt(2))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.False)

				contains, err = arr.Contains(ctx, runtime.NewString("other"))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.False)
			})

			Convey("Should return false for empty array", func() {
				arr := runtime.NewArray(0)

				contains, err := arr.Contains(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(contains, ShouldEqual, runtime.False)
			})
		})

		Convey(".IndexOf", func() {
			Convey("Should return correct index when element exists", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				index, err := arr.IndexOf(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(0))

				index, err = arr.IndexOf(ctx, runtime.NewString("test"))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(1))

				index, err = arr.IndexOf(ctx, runtime.NewFloat(3.14))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(2))
			})

			Convey("Should return -1 when element does not exist", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
				)

				index, err := arr.IndexOf(ctx, runtime.NewInt(2))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(-1))

				index, err = arr.IndexOf(ctx, runtime.NewString("other"))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(-1))
			})

			Convey("Should return -1 for empty array", func() {
				arr := runtime.NewArray(0)

				index, err := arr.IndexOf(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)
				So(index, ShouldEqual, runtime.NewInt(-1))
			})
		})

		Convey(".First", func() {
			Convey("Should return first element", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				first, err := arr.First(ctx)
				So(err, ShouldBeNil)
				So(first, ShouldEqual, runtime.NewInt(1))
			})

			Convey("Should return None for empty array", func() {
				arr := runtime.NewArray(0)

				first, err := arr.First(ctx)
				So(err, ShouldBeNil)
				So(first, ShouldEqual, runtime.None)
			})
		})

		Convey(".Last", func() {
			Convey("Should return last element", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				last, err := arr.Last(ctx)
				So(err, ShouldBeNil)
				So(last, ShouldEqual, runtime.NewFloat(3.14))
			})

			Convey("Should return None for empty array", func() {
				arr := runtime.NewArray(0)

				last, err := arr.Last(ctx)
				So(err, ShouldBeNil)
				So(last, ShouldEqual, runtime.None)
			})
		})

		Convey(".IsEmpty", func() {
			Convey("Should return true for empty array", func() {
				arr := runtime.NewArray(0)

				isEmpty, err := arr.IsEmpty(ctx)
				So(err, ShouldBeNil)
				So(isEmpty, ShouldEqual, runtime.True)
			})

			Convey("Should return false for non-empty array", func() {
				arr := runtime.NewArrayWith(runtime.NewInt(1))

				isEmpty, err := arr.IsEmpty(ctx)
				So(err, ShouldBeNil)
				So(isEmpty, ShouldEqual, runtime.False)
			})
		})

		Convey(".Clear", func() {
			Convey("Should clear all elements", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(3))

				err := arr.Clear(ctx)
				So(err, ShouldBeNil)

				length, _ = arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(0))

				isEmpty, _ := arr.IsEmpty(ctx)
				So(isEmpty, ShouldEqual, runtime.True)
			})

			Convey("Should work on empty array", func() {
				arr := runtime.NewArray(0)

				err := arr.Clear(ctx)
				So(err, ShouldBeNil)

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(0))
			})
		})

		Convey(".Copy", func() {
			Convey("Should create a shallow copy", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
				)

				copied := arr.Copy()
				copyArr := copied.(*runtime.Array)

				// Should be equal but different instances
				So(arr.Compare(copyArr), ShouldEqual, 0)

				// Modifying copy should not affect original
				copyArr.Add(ctx, runtime.NewInt(42))

				arrLength, _ := arr.Length(ctx)
				copyLength, _ := copyArr.Length(ctx)

				So(arrLength, ShouldEqual, runtime.NewInt(3))
				So(copyLength, ShouldEqual, runtime.NewInt(4))
			})

			Convey("Should copy empty array", func() {
				arr := runtime.NewArray(0)

				copied := arr.Copy()
				copyArr := copied.(*runtime.Array)

				So(arr.Compare(copyArr), ShouldEqual, 0)
				
				isEmpty, _ := copyArr.IsEmpty(ctx)
				So(isEmpty, ShouldEqual, runtime.True)
			})
		})

		Convey(".Remove", func() {
			Convey("Should remove existing element", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
					runtime.NewFloat(3.14),
					runtime.NewString("test"),
				)

				err := arr.Remove(ctx, runtime.NewString("test"))
				So(err, ShouldBeNil)

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(3))

				// Should remove first occurrence
				val0, _ := arr.Get(ctx, runtime.NewInt(0))
				So(val0, ShouldEqual, runtime.NewInt(1))

				val1, _ := arr.Get(ctx, runtime.NewInt(1))
				So(val1, ShouldEqual, runtime.NewFloat(3.14))

				val2, _ := arr.Get(ctx, runtime.NewInt(2))
				So(val2, ShouldEqual, runtime.NewString("test"))
			})

			Convey("Should return nil when element does not exist", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewString("test"),
				)

				err := arr.Remove(ctx, runtime.NewString("notfound"))
				So(err, ShouldBeNil)  // Returns nil even when not found

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(2))
			})

			Convey("Should work with empty array", func() {
				arr := runtime.NewArray(0)

				err := arr.Remove(ctx, runtime.NewInt(1))
				So(err, ShouldBeNil)  // Returns nil even when not found

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(0))
			})
		})

		Convey(".SortAsc", func() {
			Convey("Should sort integers in ascending order", func() {
				arr := runtime.NewArrayWith(
					runtime.NewInt(3),
					runtime.NewInt(1),
					runtime.NewInt(4),
					runtime.NewInt(2),
				)

				err := arr.SortAsc(ctx)
				So(err, ShouldBeNil)

				val0, _ := arr.Get(ctx, runtime.NewInt(0))
				So(val0, ShouldEqual, runtime.NewInt(1))

				val1, _ := arr.Get(ctx, runtime.NewInt(1))
				So(val1, ShouldEqual, runtime.NewInt(2))

				val2, _ := arr.Get(ctx, runtime.NewInt(2))
				So(val2, ShouldEqual, runtime.NewInt(3))

				val3, _ := arr.Get(ctx, runtime.NewInt(3))
				So(val3, ShouldEqual, runtime.NewInt(4))
			})

			Convey("Should work with empty array", func() {
				arr := runtime.NewArray(0)

				err := arr.SortAsc(ctx)
				So(err, ShouldBeNil)

				length, _ := arr.Length(ctx)
				So(length, ShouldEqual, runtime.NewInt(0))
			})
		})
	})
}