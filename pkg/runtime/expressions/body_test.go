package expressions_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBody(t *testing.T) {
	Convey("Should create a block expression", t, func() {
		s := expressions.NewBodyExpression(1)

		So(s, ShouldHaveSameTypeAs, &expressions.BodyExpression{})
	})

	Convey("Should add a new expression of a default type", t, func() {
		s := expressions.NewBodyExpression(0)

		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)
	})

	Convey("Should add a new Return expression", t, func() {
		s := expressions.NewBodyExpression(0)

		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)
	})

	Convey("Should not add an already defined Return expression", t, func() {
		s := expressions.NewBodyExpression(0)

		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeError)
		So(err.Error(), ShouldEqual, "invalid operation: return expression is already defined")
	})

	Convey("Should exec a block expression", t, func() {
		s := expressions.NewBodyExpression(1)

		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "test")
		So(err, ShouldBeNil)

		retexp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		err = s.Add(retexp)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := rootScope.Fork()
		scope.SetVariable("test", values.NewString("value"))
		fn()

		value, err := s.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
		So(value, ShouldNotBeNil)
		So(value, ShouldEqual, "value")
	})

	Convey("Should not found a missing statement", t, func() {
		s := expressions.NewBodyExpression(1)

		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := rootScope.Fork()
		scope.SetVariable("test", values.NewString("value"))
		fn()

		value, err := s.Exec(context.Background(), scope)
		So(err, ShouldNotBeNil)
		So(err.(*core.SourceErrorDetail).BaseError, ShouldHaveSameTypeAs, core.ErrNotFound)
		So(value, ShouldHaveSameTypeAs, values.None)
	})

	Convey("Should not exec a nil block expression", t, func() {
		s := expressions.NewBodyExpression(1)

		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, err := expressions.NewVariableExpression(sourceMap, "test")
		So(err, ShouldBeNil)

		err = s.Add(exp)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := rootScope.Fork()
		scope.SetVariable("test", values.NewString("value"))
		fn()

		value, err := s.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
		So(value, ShouldHaveSameTypeAs, values.None)
	})

	Convey("Should stop an execution when context is cancelled", t, func() {
		s := expressions.NewBodyExpression(1)
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
}
