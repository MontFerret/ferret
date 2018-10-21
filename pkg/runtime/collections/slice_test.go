package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func sliceIterator(value []core.Value) collections.Iterator {
	return collections.NewDefaultSliceIterator(value)
}

func TestSliceIterator(t *testing.T) {
	Convey("Should iterate over a slice", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		iter := sliceIterator(arr)

		res := make([]core.Value, 0, len(arr))

		pos := 0

		for iter.HasNext() {
			item, key, err := next(iter)

			So(err, ShouldBeNil)
			So(key.Unwrap(), ShouldEqual, pos)

			res = append(res, item)

			pos += 1
		}

		So(res, ShouldHaveLength, len(arr))
	})

	Convey("Should iterate over a slice in the same order", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		iter := sliceIterator(arr)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := next(iter)

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		for idx := range arr {
			expected := arr[idx]
			actual := res[idx]

			So(actual, ShouldEqual, expected)
		}
	})

	Convey("Should return an error when exhausted", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		iter := sliceIterator(arr)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := next(iter)

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		item, _, err := next(iter)

		So(item, ShouldBeNil)
		So(err, ShouldBeError)
	})

	Convey("Should NOT iterate over an empty slice", t, func() {
		arr := []core.Value{}

		iter := sliceIterator(arr)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}
