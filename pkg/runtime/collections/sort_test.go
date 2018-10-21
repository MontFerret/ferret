package collections_test

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func toValues(sets []collections.DataSet) []core.Value {
	res := make([]core.Value, 0, len(sets))

	for _, ds := range sets {
		res = append(res, ds.Get(collections.DefaultValueVar))
	}

	return res
}

func toArrayOfValues(sets []collections.DataSet) *values.Array {
	return values.NewArrayWith(toValues(sets)...)
}

func TestSort(t *testing.T) {
	Convey("Should sort asc", t, func() {
		arr := []core.Value{
			values.NewInt(5),
			values.NewInt(1),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(2),
		}

		s, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				return first.Get(collections.DefaultValueVar).Compare(second.Get(collections.DefaultValueVar)), nil
			},
			collections.SortDirectionAsc,
		)

		src, err := collections.NewSortIterator(
			sliceIterator(arr),
			s,
		)

		So(err, ShouldBeNil)

		res, err := collections.ToSlice(src)

		So(err, ShouldBeNil)

		numbers := []int{1, 2, 3, 4, 5}

		for idx, num := range numbers {
			So(res[idx].Get(collections.DefaultValueVar).Unwrap(), ShouldEqual, num)
		}
	})

	Convey("Should sort desc", t, func() {
		arr := []core.Value{
			values.NewInt(5),
			values.NewInt(1),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(2),
		}

		s, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				return first.Get(collections.DefaultValueVar).Compare(second.Get(collections.DefaultValueVar)), nil
			},
			collections.SortDirectionDesc,
		)

		src, err := collections.NewSortIterator(
			sliceIterator(arr),
			s,
		)

		So(err, ShouldBeNil)

		res, err := collections.ToSlice(src)

		So(err, ShouldBeNil)

		numbers := []int{5, 4, 3, 2, 1}

		for idx, num := range numbers {
			So(res[idx].Get(collections.DefaultValueVar).Unwrap(), ShouldEqual, num)
		}
	})

	Convey("Should sort asc with multiple sorters", t, func() {
		makeObj := func(one, two int) *values.Object {
			obj := values.NewObject()

			obj.Set("one", values.NewInt(one))
			obj.Set("two", values.NewInt(two))

			return obj
		}

		arr := []core.Value{
			makeObj(1, 2),
			makeObj(1, 1),
			makeObj(3, 1),
			makeObj(4, 2),
			makeObj(2, 1),
			makeObj(3, 2),
			makeObj(4, 1),
			makeObj(2, 2),
		}

		s1, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				o1, _ := first.Get(collections.DefaultValueVar).(*values.Object).Get("one")
				o2, _ := second.Get(collections.DefaultValueVar).(*values.Object).Get("one")

				return o1.Compare(o2), nil
			},
			collections.SortDirectionAsc,
		)

		s2, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				o1, _ := first.Get(collections.DefaultValueVar).(*values.Object).Get("two")
				o2, _ := second.Get(collections.DefaultValueVar).(*values.Object).Get("two")

				return o1.Compare(o2), nil
			},
			collections.SortDirectionAsc,
		)

		src, err := collections.NewSortIterator(
			sliceIterator(arr),
			s1,
			s2,
		)

		So(err, ShouldBeNil)

		sets, err := collections.ToSlice(src)

		So(err, ShouldBeNil)

		res := toValues(sets)

		j, _ := json.Marshal(res)

		So(string(j), ShouldEqual, `[{"one":1,"two":1},{"one":1,"two":2},{"one":2,"two":1},{"one":2,"two":2},{"one":3,"two":1},{"one":3,"two":2},{"one":4,"two":1},{"one":4,"two":2}]`)
	})

	Convey("Should sort desc with multiple sorters", t, func() {
		makeObj := func(one, two int) *values.Object {
			obj := values.NewObject()

			obj.Set("one", values.NewInt(one))
			obj.Set("two", values.NewInt(two))

			return obj
		}

		arr := []core.Value{
			makeObj(1, 2),
			makeObj(1, 1),
			makeObj(3, 1),
			makeObj(4, 2),
			makeObj(2, 1),
			makeObj(3, 2),
			makeObj(4, 1),
			makeObj(2, 2),
		}

		s1, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				o1, _ := first.Get(collections.DefaultValueVar).(*values.Object).Get("one")
				o2, _ := second.Get(collections.DefaultValueVar).(*values.Object).Get("one")

				return o1.Compare(o2), nil
			},
			collections.SortDirectionDesc,
		)

		s2, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				o1, _ := first.Get(collections.DefaultValueVar).(*values.Object).Get("two")
				o2, _ := second.Get(collections.DefaultValueVar).(*values.Object).Get("two")

				return o1.Compare(o2), nil
			},
			collections.SortDirectionDesc,
		)

		src, err := collections.NewSortIterator(
			sliceIterator(arr),
			s1,
			s2,
		)

		So(err, ShouldBeNil)

		sets, err := collections.ToSlice(src)

		So(err, ShouldBeNil)

		res := toValues(sets)

		j, _ := json.Marshal(res)

		So(string(j), ShouldEqual, `[{"one":4,"two":2},{"one":4,"two":1},{"one":3,"two":2},{"one":3,"two":1},{"one":2,"two":2},{"one":2,"two":1},{"one":1,"two":2},{"one":1,"two":1}]`)
	})

	Convey("Should sort asc and desc with multiple sorters", t, func() {
		makeObj := func(one, two int) *values.Object {
			obj := values.NewObject()

			obj.Set("one", values.NewInt(one))
			obj.Set("two", values.NewInt(two))

			return obj
		}

		arr := []core.Value{
			makeObj(1, 2),
			makeObj(1, 1),
			makeObj(3, 1),
			makeObj(4, 2),
			makeObj(2, 1),
			makeObj(3, 2),
			makeObj(4, 1),
			makeObj(2, 2),
		}

		s1, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				o1, _ := first.Get(collections.DefaultValueVar).(*values.Object).Get("one")
				o2, _ := second.Get(collections.DefaultValueVar).(*values.Object).Get("one")

				return o1.Compare(o2), nil
			},
			collections.SortDirectionAsc,
		)

		s2, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				o1, _ := first.Get(collections.DefaultValueVar).(*values.Object).Get("two")
				o2, _ := second.Get(collections.DefaultValueVar).(*values.Object).Get("two")

				return o1.Compare(o2), nil
			},
			collections.SortDirectionDesc,
		)

		src, err := collections.NewSortIterator(
			sliceIterator(arr),
			s1,
			s2,
		)

		So(err, ShouldBeNil)

		sets, err := collections.ToSlice(src)

		So(err, ShouldBeNil)

		res := toValues(sets)

		j, _ := json.Marshal(res)

		So(string(j), ShouldEqual, `[{"one":1,"two":2},{"one":1,"two":1},{"one":2,"two":2},{"one":2,"two":1},{"one":3,"two":2},{"one":3,"two":1},{"one":4,"two":2},{"one":4,"two":1}]`)
	})

	Convey("Should sort desc and asc with multiple sorters", t, func() {
		makeObj := func(one, two int) *values.Object {
			obj := values.NewObject()

			obj.Set("one", values.NewInt(one))
			obj.Set("two", values.NewInt(two))

			return obj
		}

		arr := []core.Value{
			makeObj(1, 2),
			makeObj(1, 1),
			makeObj(3, 1),
			makeObj(4, 2),
			makeObj(2, 1),
			makeObj(3, 2),
			makeObj(4, 1),
			makeObj(2, 2),
		}

		s1, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				o1, _ := first.Get(collections.DefaultValueVar).(*values.Object).Get("one")
				o2, _ := second.Get(collections.DefaultValueVar).(*values.Object).Get("one")

				return o1.Compare(o2), nil
			},
			collections.SortDirectionDesc,
		)

		s2, _ := collections.NewSorter(
			func(first collections.DataSet, second collections.DataSet) (int, error) {
				o1, _ := first.Get(collections.DefaultValueVar).(*values.Object).Get("two")
				o2, _ := second.Get(collections.DefaultValueVar).(*values.Object).Get("two")

				return o1.Compare(o2), nil
			},
			collections.SortDirectionAsc,
		)

		src, err := collections.NewSortIterator(
			sliceIterator(arr),
			s1,
			s2,
		)

		So(err, ShouldBeNil)

		sets, err := collections.ToSlice(src)

		So(err, ShouldBeNil)

		res := toValues(sets)

		j, _ := json.Marshal(res)

		So(string(j), ShouldEqual, `[{"one":4,"two":1},{"one":4,"two":2},{"one":3,"two":1},{"one":3,"two":2},{"one":2,"two":1},{"one":2,"two":2},{"one":1,"two":1},{"one":1,"two":2}]`)
	})
}
