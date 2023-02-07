package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func objectIterator(obj *values.Object) collections.Iterator {
	iter, _ := collections.NewDefaultKeyedIterator(obj)

	return iter
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

		iter := objectIterator(m)

		res := make([]core.Value, 0, m.Length())

		ctx := context.Background()
		scope, _ := core.NewRootScope()

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

			expected, exists := m.Get(values.NewString(key.String()))

			So(bool(exists), ShouldBeTrue)
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

		iter := objectIterator(m)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		nextScope, err := iter.Next(ctx, scope)

		So(nextScope, ShouldBeNil)
		So(err, ShouldEqual, core.ErrNoMoreData)
	})

	Convey("Should NOT iterate over a empty map", t, func() {
		m := values.NewObject()

		iter := objectIterator(m)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		nextScope, err := iter.Next(ctx, scope)

		So(nextScope, ShouldBeNil)
		So(err, ShouldEqual, core.ErrNoMoreData)
	})
}
