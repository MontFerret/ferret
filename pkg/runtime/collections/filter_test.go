package collections_test

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

func TestFilter(t *testing.T) {
	Convey("Should filter out non-even values", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(val core.Value, _ core.Value) (bool, error) {
			i := float64(val.Unwrap().(int))
			calc := float64(i / 2)

			return calc == math.Floor(calc), nil
		}

		iter, err := collections.NewFilterIterator(
			collections.NewSliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(res, ShouldHaveLength, 2)
	})

	Convey("Should filter out non-even keys", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(_ core.Value, key core.Value) (bool, error) {
			i := float64(key.Unwrap().(int))

			if i == 0 {
				return false, nil
			}

			calc := float64(i / 2)

			return calc == math.Floor(calc), nil
		}

		iter, err := collections.NewFilterIterator(
			collections.NewSliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(res, ShouldHaveLength, 2)
	})

	Convey("Should filter out values all values", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(val core.Value, _ core.Value) (bool, error) {
			return false, nil
		}

		iter, err := collections.NewFilterIterator(
			collections.NewSliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(res, ShouldHaveLength, 0)
	})

	Convey("Should pass through all values", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(val core.Value, _ core.Value) (bool, error) {
			return true, nil
		}

		iter, err := collections.NewFilterIterator(
			collections.NewSliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		So(res, ShouldHaveLength, len(arr))
	})

	Convey("Should return an error when exhausted", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(val core.Value, _ core.Value) (bool, error) {
			return true, nil
		}

		iter, err := collections.NewFilterIterator(
			collections.NewSliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		res := make([]core.Value, 0, len(arr))

		for iter.HasNext() {
			item, _, err := iter.Next()

			So(err, ShouldBeNil)

			res = append(res, item)
		}

		_, _, err = iter.Next()

		So(err, ShouldBeError)
	})

	Convey("Should iterate over nested filter", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		// i < 5
		predicate1 := func(val core.Value, _ core.Value) (bool, error) {
			return val.Compare(values.NewInt(5)) == -1, nil
		}

		// i > 2
		predicate2 := func(val core.Value, _ core.Value) (bool, error) {
			return val.Compare(values.NewInt(2)) == 1, nil
		}

		it, _ := collections.NewFilterIterator(
			collections.NewSliceIterator(arr),
			predicate1,
		)

		iter, err := collections.NewFilterIterator(
			it,
			predicate2,
		)

		So(err, ShouldBeNil)

		res, err := collections.ToSlice(iter)

		So(err, ShouldBeNil)

		js, _ := json.Marshal(res)

		So(string(js), ShouldEqual, `[3,4]`)
	})
}
