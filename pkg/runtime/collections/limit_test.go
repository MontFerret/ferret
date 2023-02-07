package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestLimit(t *testing.T) {
	Convey("Should limit iteration", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		iter, err := collections.NewLimitIterator(
			sliceIterator(arr),
			1,
			0,
		)

		So(err, ShouldBeNil)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 1)
	})

	Convey("Should limit iteration (2)", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		iter, err := collections.NewLimitIterator(
			sliceIterator(arr),
			2,
			0,
		)

		So(err, ShouldBeNil)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)
	})

	Convey("Should limit iteration with offset", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		offset := 2
		iter, err := collections.NewLimitIterator(
			sliceIterator(arr),
			2,
			offset,
		)

		So(err, ShouldBeNil)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)

		for idx, nextScope := range res {
			expected := arr[idx+offset]
			current := nextScope.MustGetVariable(collections.DefaultValueVar)

			So(expected, ShouldEqual, current)
		}
	})

	Convey("Should limit iteration with offset at the end", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		offset := 3

		iter, err := collections.NewLimitIterator(
			sliceIterator(arr),
			2,
			offset,
		)

		So(err, ShouldBeNil)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)

		for idx, nextScope := range res {
			expected := arr[idx+offset]
			current := nextScope.MustGetVariable(collections.DefaultValueVar)

			So(expected, ShouldEqual, current)
		}
	})

	Convey("Should limit iteration with offset with going out of bounds", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		}

		offset := 4

		iter, err := collections.NewLimitIterator(
			sliceIterator(arr),
			2,
			offset,
		)

		So(err, ShouldBeNil)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 1)
	})
}
