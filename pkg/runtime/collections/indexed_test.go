package collections_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
			ds, err := iter.Next(ctx, scope)

			So(err, ShouldBeNil)

			if ds == nil {
				break
			}

			res = append(res, ds.Get(collections.DefaultValueVar))

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

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		for {
			ds, err := iter.Next(ctx, scope)

			So(err, ShouldBeNil)

			if ds == nil {
				break
			}

			res = append(res, ds.Get(collections.DefaultValueVar))
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
			ds, err := iter.Next(ctx, scope)

			So(err, ShouldBeNil)

			if ds == nil {
				break
			}

			res = append(res, ds.Get(collections.DefaultValueVar))
		}

		item, err := iter.Next(ctx, scope)

		So(item, ShouldBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Should NOT iterate over an empty array", t, func() {
		arr := values.NewArray(10)

		iter := arrayIterator(arr)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		ds, err := iter.Next(ctx, scope)

		So(err, ShouldBeNil)
		So(ds, ShouldBeNil)
	})
}
