package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func mapIterator(m map[string]core.Value) collections.Iterator {
	iter, _ := collections.NewDefaultMapIterator(m)

	return iter
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

		iter := mapIterator(m)

		res := make([]core.Value, 0, len(m))
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

		iter := mapIterator(m)
		ctx := context.Background()
		scope, _ := core.NewRootScope()

		_, err := collections.ToSlice(ctx, scope, iter)
		So(err, ShouldBeNil)

		item, err := iter.Next(ctx, scope)

		So(item, ShouldBeNil)
		So(err, ShouldEqual, core.ErrNoMoreData)
	})

	Convey("Should NOT iterate over a empty map", t, func() {
		m := make(map[string]core.Value)

		iter := mapIterator(m)

		ctx := context.Background()
		scope, _ := core.NewRootScope()

		item, err := iter.Next(ctx, scope)

		So(item, ShouldBeNil)
		So(err, ShouldEqual, core.ErrNoMoreData)
	})
}
