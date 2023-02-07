package expressions_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/collections"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type IterableFn func(ctx context.Context, scope *core.Scope) (collections.Iterator, error)

func (f IterableFn) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	return f(ctx, scope)
}

type ExpressionFn func(ctx context.Context, scope *core.Scope) (core.Value, error)

func (f ExpressionFn) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	return f(ctx, scope)
}

func TestBlock(t *testing.T) {
	newExp := func(values []core.Value) (*expressions.BlockExpression, error) {
		iter, err := collections.NewDefaultSliceIterator(values)

		if err != nil {
			return nil, err
		}

		return expressions.NewBlockExpression(IterableFn(func(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
			return iter, nil
		}))
	}

	Convey("Should create a block expression", t, func() {
		s, err := newExp(make([]core.Value, 0, 10))

		So(err, ShouldBeNil)
		So(s, ShouldHaveSameTypeAs, &expressions.BlockExpression{})
	})

	Convey("Should add a new expression of a default type", t, func() {
		s, _ := newExp(make([]core.Value, 0, 10))

		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		s.Add(exp)
	})

	Convey("Should exec a block expression", t, func() {
		s, _ := newExp(make([]core.Value, 0, 10))

		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, err := expressions.NewVariableDeclarationExpression(sourceMap, "test", ExpressionFn(func(ctx context.Context, scope *core.Scope) (core.Value, error) {
			return values.NewString("value"), nil
		}))
		So(err, ShouldBeNil)

		s.Add(exp)

		rootScope, _ := core.NewRootScope()
		scope := rootScope.Fork()

		_, err = s.Exec(context.Background(), scope)
		So(err, ShouldBeNil)

		val, err := scope.GetVariable("test")
		So(err, ShouldBeNil)

		So(val, ShouldEqual, "value")
	})

	Convey("Should not exec a nil block expression", t, func() {
		s, _ := newExp(make([]core.Value, 0, 10))

		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, err := expressions.NewVariableExpression(sourceMap, "test")
		So(err, ShouldBeNil)

		s.Add(exp)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := rootScope.Fork()
		scope.SetVariable("test", values.NewString("value"))
		fn()

		value, err := s.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
		So(value, ShouldHaveSameTypeAs, values.None)
	})

	Convey("Should return an iterator", t, func() {
		s, _ := newExp([]core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
		})
		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, _ := expressions.NewVariableExpression(sourceMap, "test")
		s.Add(exp)

		rootScope, _ := core.NewRootScope()
		scope := rootScope.Fork()
		scope.SetVariable("test", values.NewString("value"))

		iter, err := s.Iterate(context.Background(), scope)
		So(err, ShouldBeNil)

		items, err := collections.ToSlice(context.Background(), scope, iter)
		So(err, ShouldBeNil)
		So(items, ShouldHaveLength, 3)
	})

	Convey("Should stop an execution when context is cancelled", t, func() {
		s, _ := newExp(make([]core.Value, 0, 10))
		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, _ := expressions.NewVariableExpression(sourceMap, "test")
		s.Add(exp)

		rootScope, _ := core.NewRootScope()
		scope := rootScope.Fork()
		scope.SetVariable("test", values.NewString("value"))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := s.Exec(ctx, scope)
		So(err, ShouldEqual, core.ErrTerminated)
	})

	Convey("Should stop an execution when context is cancelled 2", t, func() {
		s, _ := newExp([]core.Value{
			values.NewInt(1),
			values.NewInt(2),
			values.NewInt(3),
		})
		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, _ := expressions.NewVariableExpression(sourceMap, "test")
		s.Add(exp)

		rootScope, _ := core.NewRootScope()
		scope := rootScope.Fork()
		scope.SetVariable("test", values.NewString("value"))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := s.Iterate(ctx, scope)
		So(err, ShouldEqual, core.ErrTerminated)
	})
}
