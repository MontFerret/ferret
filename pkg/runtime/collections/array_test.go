package collections_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func next(iterator collections.Iterator) (core.Value, core.Value, error) {
	ds, err := iterator.Next()

	if err != nil {
		return nil, nil, err
	}

	val := ds[collections.DefaultValueVar]
	key := ds[collections.DefaultKeyVar]

	return val, key, nil
}

func arrayIterator(arr *values.Array) collections.Iterator {
	return collections.NewDefaultArrayIterator(arr)
}

func TestArrayIterator(t *testing.T) {

	Convey("Should iterate over an array", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		iter := arrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

		pos := 0

		for iter.HasNext() {
			item, key, err := next(iter)

			So(err, ShouldBeNil)
			So(key.Unwrap(), ShouldEqual, pos)

			res = append(res, item)

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

		iter := arrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

		for iter.HasNext() {
			item, _, err := next(iter)

			So(err, ShouldBeNil)

			res = append(res, item)
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

		iter := arrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

		for iter.HasNext() {
			item, _, err := next(iter)

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		item, _, err := next(iter)

		So(item, ShouldBeNil)
		So(err, ShouldBeError)
	})

	Convey("Should NOT iterate over an empty array", t, func() {
		arr := values.NewArray(10)

		iter := arrayIterator(arr)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}
