package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestUniqueIterator(t *testing.T) {
	Convey("Should return only unique items", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(3),
			values.NewInt(5),
			values.NewInt(6),
			values.NewInt(5),
			values.NewInt(6),
		}

		iter, err := collections.NewUniqueIterator(
			sliceIterator(arr),
			collections.DefaultValueVar,
		)

		So(err, ShouldBeNil)
		ctx := context.Background()
		scope, _ := core.NewRootScope()

		sets, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)

		res := toArrayOfValues(sets)

		So(res.String(), ShouldEqual, `[1,2,3,4,5,6]`)
	})

	Convey("Should return only unique items 2", t, func() {
		arr := []core.Value{
			values.NewInt(1),
			values.NewInt(1),
			values.NewInt(1),
			values.NewInt(1),
			values.NewInt(1),
			values.NewInt(1),
		}

		iter, err := collections.NewUniqueIterator(
			sliceIterator(arr),
			collections.DefaultValueVar,
		)

		So(err, ShouldBeNil)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		sets, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)

		res := toArrayOfValues(sets)

		So(res.String(), ShouldEqual, `[1]`)
	})

	Convey("Should return only unique items 3", t, func() {
		arr := []core.Value{
			values.NewString("a"),
			values.NewString("b"),
			values.NewString("c"),
			values.NewString("d"),
			values.NewString("e"),
			values.NewString("a"),
			values.NewString("b"),
			values.NewString("f"),
			values.NewString("d"),
			values.NewString("e"),
			values.NewString("f"),
		}

		iter, err := collections.NewUniqueIterator(
			sliceIterator(arr),
			collections.DefaultValueVar,
		)

		So(err, ShouldBeNil)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		sets, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)

		res := toArrayOfValues(sets)

		So(res.String(), ShouldEqual, `["a","b","c","d","e","f"]`)
	})
}
