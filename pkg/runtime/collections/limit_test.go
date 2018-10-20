package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLimit(t *testing.T) {
	Convey("Should limit iteration", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		src, err := collections.NewLimitIterator(
			collections.NewSliceIterator(arr),
			1,
			0,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for src.HasNext() {
			item, _, err := src.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(len(res), ShouldEqual, 1)
	})

	Convey("Should limit iteration (2)", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		src, err := collections.NewLimitIterator(
			collections.NewSliceIterator(arr),
			2,
			0,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for src.HasNext() {
			item, _, err := src.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(len(res), ShouldEqual, 2)
	})

	Convey("Should limit iteration with offset", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		offset := 2
		src, err := collections.NewLimitIterator(
			collections.NewSliceIterator(arr),
			2,
			offset,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for src.HasNext() {
			item, _, err := src.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(len(res), ShouldEqual, 2)

		for idx, current := range res {
			expected := arr[idx+offset]

			So(expected, ShouldEqual, current)
		}
	})

	Convey("Should limit iteration with offset at the end", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		offset := 3

		src, err := collections.NewLimitIterator(
			collections.NewSliceIterator(arr),
			2,
			offset,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for src.HasNext() {
			item, _, err := src.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(len(res), ShouldEqual, 2)

		for idx, current := range res {
			expected := arr[idx+offset]

			So(expected, ShouldEqual, current)
		}
	})

	Convey("Should limit iteration with offset with going out of bounds", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		offset := 4

		src, err := collections.NewLimitIterator(
			collections.NewSliceIterator(arr),
			2,
			offset,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for src.HasNext() {
			item, _, err := src.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(len(res), ShouldEqual, 1)
	})
}
