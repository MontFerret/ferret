package collections_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDoWhileIterator(t *testing.T) {
	Convey("Should iterate while a predicate returns true", t, func() {
		predicateCounter := 0
		iterationCounter := 0
		expectedCount := 10

		iterator, err := collections.NewDefaultWhileIterator(collections.WhileModePre, func(ctx context.Context, scope *core.Scope) (bool, error) {
			if predicateCounter == (expectedCount - 1) {
				return false, nil
			}

			predicateCounter++

			return true, nil
		})

		So(err, ShouldBeNil)

		scope, fn := core.NewRootScope()
		defer fn()
		err = collections.ForEach(context.Background(), scope, iterator, func(ctx context.Context, scope *core.Scope) bool {
			iterationCounter++

			return true
		})

		So(err, ShouldBeNil)
		So(iterationCounter, ShouldEqual, expectedCount)
	})

	Convey("Should iterate once if a predicate returns false", t, func() {
		iterationCounter := 0
		expectedCount := 1

		iterator, err := collections.NewDefaultWhileIterator(collections.WhileModePre, func(ctx context.Context, scope *core.Scope) (bool, error) {
			return false, nil
		})

		So(err, ShouldBeNil)

		scope, fn := core.NewRootScope()
		defer fn()
		err = collections.ForEach(context.Background(), scope, iterator, func(ctx context.Context, scope *core.Scope) bool {
			iterationCounter++

			return true
		})

		So(err, ShouldBeNil)
		So(iterationCounter, ShouldEqual, expectedCount)
	})
}
