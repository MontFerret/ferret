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
			item, key, err := iter.Next()

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

		iter := collections.NewSliceIterator(arr)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := iter.Next()

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

		iter := collections.NewSliceIterator(arr)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		item, _, err := iter.Next()

		So(item, ShouldEqual, values.None)
		So(err, ShouldBeError)
	})

	Convey("Should NOT iterate over an empty slice", t, func() {
		arr := []core.Value{}

		iter := collections.NewSliceIterator(arr)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}

func TestMapIterator(t *testing.T) {
	Convey("Should iterate over a map", t, func() {
		m := map[string]core.Value{
			"one":   values.NewInt(1),
			"two":   values.NewInt(2),
			"three": values.NewInt(3),
			"four":  values.NewInt(4),
			"five":  values.NewInt(5),
		}

		iter := collections.NewMapIterator(m)

		res := make([]core.Value, 0, len(m))

		for iter.HasNext() {
			item, key, err := iter.Next()

			So(err, ShouldBeNil)

			expected, exists := m[key.String()]

			So(exists, ShouldBeTrue)
			So(expected, ShouldEqual, item)

			res = append(res, item)
		}

		So(res, ShouldHaveLength, len(m))
	})

	Convey("Should return an error when exhausted", t, func() {
		m := map[string]core.Value{
			"one":   values.NewInt(1),
			"two":   values.NewInt(2),
			"three": values.NewInt(3),
			"four":  values.NewInt(4),
			"five":  values.NewInt(5),
		}

		iter := collections.NewMapIterator(m)

		res := make([]core.Value, 0, len(m))

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		item, _, err := iter.Next()

		So(item, ShouldEqual, values.None)
		So(err, ShouldBeError)
	})

	Convey("Should NOT iterate over a empty map", t, func() {
		m := make(map[string]core.Value)

		iter := collections.NewMapIterator(m)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
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

		iter := collections.NewArrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

		pos := 0

		for iter.HasNext() {
			item, key, err := iter.Next()

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

		iter := collections.NewArrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

		for iter.HasNext() {
			item, _, err := iter.Next()

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

		iter := collections.NewArrayIterator(arr)

		res := make([]core.Value, 0, arr.Length())

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		item, _, err := iter.Next()

		So(item, ShouldEqual, values.None)
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

func TestObjectIterator(t *testing.T) {
	Convey("Should iterate over a map", t, func() {
		m := values.NewObjectWith(
			values.NewObjectProperty("one", values.NewInt(1)),
			values.NewObjectProperty("two", values.NewInt(2)),
			values.NewObjectProperty("three", values.NewInt(3)),
			values.NewObjectProperty("four", values.NewInt(4)),
			values.NewObjectProperty("five", values.NewInt(5)),
		)

		iter := collections.NewObjectIterator(m)

		res := make([]core.Value, 0, m.Length())

		for iter.HasNext() {
			item, key, err := iter.Next()

			So(err, ShouldBeNil)

			expected, exists := m.Get(values.NewString(key.String()))

			So(exists, ShouldBeTrue)
			So(expected, ShouldEqual, item)

			res = append(res, item)
		}

		So(res, ShouldHaveLength, m.Length())
	})

	Convey("Should return an error when exhausted", t, func() {
		m := values.NewObjectWith(
			values.NewObjectProperty("one", values.NewInt(1)),
			values.NewObjectProperty("two", values.NewInt(2)),
			values.NewObjectProperty("three", values.NewInt(3)),
			values.NewObjectProperty("four", values.NewInt(4)),
			values.NewObjectProperty("five", values.NewInt(5)),
		)

		iter := collections.NewObjectIterator(m)

		res := make([]core.Value, 0, m.Length())

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		item, _, err := iter.Next()

		So(item, ShouldEqual, values.None)
		So(err, ShouldBeError)
	})

	Convey("Should NOT iterate over a empty map", t, func() {
		m := values.NewObject()

		iter := collections.NewObjectIterator(m)

		var iterated bool

		for iter.HasNext() {
			iterated = true
		}

		So(iterated, ShouldBeFalse)
	})
}
