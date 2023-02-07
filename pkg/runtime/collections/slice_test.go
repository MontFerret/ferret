package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func sliceIterator(value []core.Value) collections.Iterator {
	iter, _ := collections.NewDefaultSliceIterator(value)

	return iter
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
		ctx := context.Background()
		scope, _ := core.NewRootScope()

		pos := 0

		for {
			nextScope, err := iter.Next(ctx, scope)

			if err != nil {
				if core.IsNoMoreData(err) {
					break
				}

				So(err, ShouldBeNil)
			}

			key := nextScope.MustGetVariable(collections.DefaultKeyVar)
			item := nextScope.MustGetVariable(collections.DefaultValueVar)

			So(key.Unwrap(), ShouldEqual, pos)

			res = append(res, item)

			pos++
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
		ctx := context.Background()
		scope, _ := core.NewRootScope()

		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)

		for idx := range arr {
			expected := arr[idx]
			nextScope := res[idx]
			actual := nextScope.MustGetVariable(collections.DefaultValueVar)

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
		ctx := context.Background()
		scope, _ := core.NewRootScope()

		_, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)

		item, err := iter.Next(ctx, scope)

		So(item, ShouldBeNil)
		So(err, ShouldEqual, core.ErrNoMoreData)
	})

	Convey("Should NOT iterate over an empty slice", t, func() {
		arr := []core.Value{}

		iter := sliceIterator(arr)
		ctx := context.Background()
		scope, _ := core.NewRootScope()

		item, err := iter.Next(ctx, scope)

		So(item, ShouldBeNil)
		So(err, ShouldEqual, core.ErrNoMoreData)
	})
}
