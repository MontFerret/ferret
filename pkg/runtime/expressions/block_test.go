package expressions_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBlockExpression(t *testing.T) {
	Convey("Should match", t, func() {
		s := expressions.NewBlockExpression(1)

		So(s, ShouldHaveSameTypeAs, &expressions.BlockExpression{})
	})
}

func TestNewBlockExpressionWith(t *testing.T) {
	Convey("Should match", t, func() {
		s := expressions.NewBlockExpressionWith(&expressions.BlockExpression{})

		So(s, ShouldHaveSameTypeAs, &expressions.BlockExpression{})
	})
}

func TestBlockExpressionAddVariableExpression(t *testing.T) {
	Convey("Should match", t, func() {
		s := expressions.NewBlockExpression(0)

		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)
	})
}

func TestBlockExpressionAddReturnExpression(t *testing.T) {
	Convey("Should match", t, func() {
		s := expressions.NewBlockExpression(0)

		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)
	})
}

func TestBlockExpressionAddReturnExpressionFailed(t *testing.T) {
	Convey("Should match", t, func() {
		s := expressions.NewBlockExpression(0)

		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeError)
		So(err.Error(), ShouldEqual, "return expression is already defined: invalid operation")
	})
}

func TestBlockExpressionExec(t *testing.T) {
	Convey("Should match", t, func() {
		s := expressions.NewBlockExpression(1)

		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "test")
		So(err, ShouldBeNil)

		retexp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		err = s.Add(retexp)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := core.NewScope(rootScope)
		scope.SetVariable("test", values.NewString("value"))
		fn()

		value, err := s.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
		So(value, ShouldNotBeNil)
		So(value, ShouldEqual, "value")
	})
}

func TestBlockExpressionExecNonFound(t *testing.T) {
	Convey("Should match", t, func() {
		s := expressions.NewBlockExpression(1)

		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := core.NewScope(rootScope)
		scope.SetVariable("test", values.NewString("value"))
		fn()

		value, err := s.Exec(context.Background(), scope)
		So(err, ShouldNotBeNil)
		So(err, ShouldHaveSameTypeAs, core.ErrNotFound)
		So(value, ShouldHaveSameTypeAs, values.None)
	})
}

func TestBlockExpressionExecNilExpression(t *testing.T) {
	Convey("Should match", t, func() {
		s := expressions.NewBlockExpression(1)

		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, err := expressions.NewVariableExpression(sourceMap, "test")
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := core.NewScope(rootScope)
		scope.SetVariable("test", values.NewString("value"))
		fn()

		value, err := s.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
		So(value, ShouldHaveSameTypeAs, values.None)
	})
}
