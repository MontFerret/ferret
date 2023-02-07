package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func arrayIterator(arr *values.Array) collections.Iterator {
	iterator, _ := collections.NewDefaultIndexedIterator(arr)

	return iterator
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

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		for {
			nextScope, err := iter.Next(ctx, scope.Fork())

			if err != nil {
				if core.IsNoMoreData(err) {
					break
				}

				So(err, ShouldBeNil)
			}

			res = append(res, nextScope.MustGetVariable(collections.DefaultValueVar))

			pos++
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

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		for {
			nextScope, err := iter.Next(ctx, scope.Fork())

			if err != nil {
				if core.IsNoMoreData(err) {
					break
				}

				So(err, ShouldBeNil)
			}

			res = append(res, nextScope.MustGetVariable(collections.DefaultValueVar))
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

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		for {
			nextScope, err := iter.Next(ctx, scope.Fork())

			if err != nil {
				if core.IsNoMoreData(err) {
					break
				}

				So(err, ShouldBeNil)
			}

			res = append(res, nextScope.MustGetVariable(collections.DefaultValueVar))
		}

		item, err := iter.Next(ctx, scope)

		So(item, ShouldBeNil)
		So(err, ShouldEqual, core.ErrNoMoreData)
		So(res, ShouldHaveLength, int(arr.Length()))
	})

	Convey("Should NOT iterate over an empty array", t, func() {
		arr := values.NewArray(10)

		iter := arrayIterator(arr)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		nextScope, err := iter.Next(ctx, scope)

		So(err, ShouldEqual, core.ErrNoMoreData)
		So(nextScope, ShouldBeNil)
	})
}
