package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestArrayIterator(t *testing.T) {
	Convey("Should iterate over an array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		iter := collections.NewArrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

		pos := 0

		for iter.HasNext() {
			set, err := iter.Next()

			So(err, ShouldBeNil)
			So(set, ShouldHaveLength, 2)
			So(set[1].Unwrap(), ShouldEqual, pos)

			res = append(res, set[0])

			pos += 1
		}

		So(res, ShouldHaveLength, arr.Length())
	})

	Convey("Should iterate over an array in the same order", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		iter := collections.NewArrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

		for iter.HasNext() {
			set, err := iter.Next()

			So(err, ShouldBeNil)
			So(set, ShouldHaveLength, 2)

			res = append(res, set[0])
		}

		arr.ForEach(func(expected core.Value, idx int) bool {
			actual := res[idx]

			So(actual, ShouldEqual, expected)

			return true
		})
	})

	Convey("Should return an error when exhausted", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		iter := collections.NewArrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

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

	Convey("Should NOT iterate over an empty array", t, func() {
		arr := values.NewArray(10)

		iter := collections.NewArrayIterator(arr)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}
