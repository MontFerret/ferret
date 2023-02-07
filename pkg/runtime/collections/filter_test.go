package collections_test

import (
	"context"
	"encoding/json"
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestFilter(t *testing.T) {
	Convey("Should filter out non-even result", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(_ context.Context, scope *core.Scope) (bool, error) {
			i := float64(scope.MustGetVariable(collections.DefaultValueVar).Unwrap().(int))
			calc := i / 2

			return calc == math.Floor(calc), nil
		}

		iter, err := collections.NewFilterIterator(
			sliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		res, err := collections.ToSlice(context.Background(), scope, iter)

		So(err, ShouldBeNil)
		So(res, ShouldHaveLength, 2)
	})

	Convey("Should filter out non-even groupKeys", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(_ context.Context, scope *core.Scope) (bool, error) {
			i := float64(scope.MustGetVariable(collections.DefaultKeyVar).Unwrap().(int))

			if i == 0 {
				return false, nil
			}

			calc := i / 2

			return calc == math.Floor(calc), nil
		}

		iter, err := collections.NewFilterIterator(
			sliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		res, err := collections.ToSlice(context.Background(), scope, iter)

		So(err, ShouldBeNil)

		So(res, ShouldHaveLength, 2)
	})

	Convey("Should filter out result all result", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(_ context.Context, _ *core.Scope) (bool, error) {
			return false, nil
		}

		iter, err := collections.NewFilterIterator(
			sliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		res, err := collections.ToSlice(context.Background(), scope, iter)

		So(err, ShouldBeNil)
		So(res, ShouldHaveLength, 0)
	})

	Convey("Should pass through all result", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		predicate := func(_ context.Context, _ *core.Scope) (bool, error) {
			return true, nil
		}

		iter, err := collections.NewFilterIterator(
			sliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		res, err := collections.ToSlice(context.Background(), scope, iter)

		So(err, ShouldBeNil)
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

		predicate := func(_ context.Context, _ *core.Scope) (bool, error) {
			return true, nil
		}

		iter, err := collections.NewFilterIterator(
			sliceIterator(arr),
			predicate,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		_, err = collections.ToSlice(context.Background(), scope, iter)

		So(err, ShouldBeNil)

		item, err := iter.Next(context.Background(), scope)

		So(item, ShouldBeNil)
		So(err, ShouldEqual, core.ErrNoMoreData)
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
		predicate1 := func(_ context.Context, scope *core.Scope) (bool, error) {
			return scope.MustGetVariable(collections.DefaultValueVar).Compare(values.NewInt(5)) == -1, nil
		}

		// i > 2
		predicate2 := func(_ context.Context, scope *core.Scope) (bool, error) {
			return scope.MustGetVariable(collections.DefaultValueVar).Compare(values.NewInt(2)) == 1, nil
		}

		it, _ := collections.NewFilterIterator(
			sliceIterator(arr),
			predicate1,
		)

		iter, err := collections.NewFilterIterator(
			it,
			predicate2,
		)

		So(err, ShouldBeNil)

		scope, _ := core.NewRootScope()
		sets, err := collections.ToSlice(context.Background(), scope, iter)

		So(err, ShouldBeNil)

		res := toArrayOfValues(sets)

		js, _ := json.Marshal(res)

		So(string(js), ShouldEqual, `[3,4]`)
	})
}
