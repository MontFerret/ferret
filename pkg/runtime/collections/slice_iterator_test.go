package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSliceIterator(t *testing.T) {
	Convey("Should iterate over a slice", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		iter := collections.NewSliceIterator(arr)

		res := make([]core.Value, 0, len(arr))

		pos := 0

		for iter.HasNext() {
			set, err := iter.Next()

			So(err, ShouldBeNil)
			So(set, ShouldHaveLength, 2)
			So(set[1].Unwrap(), ShouldEqual, pos)

			res = append(res, set[0])

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

		iter := collections.NewSliceIterator(arr)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			set, err := iter.Next()

			So(err, ShouldBeNil)
			So(set, ShouldHaveLength, 2)

			res = append(res, set[0])
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

		iter := collections.NewSliceIterator(arr)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			set, err := iter.Next()

			So(err, ShouldBeNil)
			So(set, ShouldHaveLength, 2)

			res = append(res, set[0])
		}

		set, err := iter.Next()

		So(set, ShouldBeNil)
		So(err, ShouldBeError)
	})

	Convey("Should NOT iterate over an empty slice", t, func() {
		arr := make([]core.Value, 0, 0)

		iter := collections.NewSliceIterator(arr)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}
