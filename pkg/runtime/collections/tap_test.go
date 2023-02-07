package collections_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func tapIterator(values collections.Iterator, predicate core.Expression) collections.Iterator {
	iter, _ := collections.NewTapIterator(values, predicate)

	return iter
}

type TestExpression struct {
	fn func(ctx context.Context, scope *core.Scope) (core.Value, error)
}

func (exp *TestExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	return exp.fn(ctx, scope)
}

type ErrorIterator struct{}

func (iterator *ErrorIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	return nil, core.ErrInvalidOperation
}

func TestTapIterator(t *testing.T) {
	Convey("Should iterate over a given iterator and execute a predicate", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		counter := 0

		iter := tapIterator(arrayIterator(arr), &TestExpression{
			fn: func(ctx context.Context, scope *core.Scope) (core.Value, error) {
				counter++

				return values.None, nil
			},
		})

		ctx := context.Background()
		scope, _ := core.NewRootScope()
		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)
		So(res, ShouldHaveLength, int(arr.Length()))
		So(counter, ShouldEqual, int(arr.Length()))
	})

	Convey("Should stop when a predicate return an error", t, func() {
		arr := values.NewArrayWith(
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
			values.NewInt(4),
			values.NewInt(5),
		)

		counter := 0

		iter := tapIterator(arrayIterator(arr), &TestExpression{
			fn: func(ctx context.Context, scope *core.Scope) (core.Value, error) {
				counter++

				return values.None, core.ErrInvalidOperation
			},
		})

		ctx := context.Background()
		scope, _ := core.NewRootScope()
		_, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldNotBeNil)
		So(counter, ShouldEqual, 1)
	})

	Convey("Should not invoke a predicate when underlying iterator returns error", t, func() {
		counter := 0

		iter := tapIterator(&ErrorIterator{}, &TestExpression{
			fn: func(ctx context.Context, scope *core.Scope) (core.Value, error) {
				counter++

				return values.None, core.ErrInvalidOperation
			},
		})

		ctx := context.Background()
		scope, _ := core.NewRootScope()
		_, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldNotBeNil)
		So(counter, ShouldEqual, 0)
	})

	Convey("Should not invoke a predicate when underlying iterator is empty", t, func() {
		arr := values.NewArray(0)

		counter := 0

		iter := tapIterator(arrayIterator(arr), &TestExpression{
			fn: func(ctx context.Context, scope *core.Scope) (core.Value, error) {
				counter++

				return values.None, core.ErrInvalidOperation
			},
		})

		ctx := context.Background()
		scope, _ := core.NewRootScope()
		res, err := collections.ToSlice(ctx, scope, iter)

		So(err, ShouldBeNil)
		So(res, ShouldHaveLength, 0)
		So(counter, ShouldEqual, 0)
	})
}
