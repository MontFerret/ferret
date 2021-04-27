package expressions_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReturnExpression(t *testing.T) {
	Convey("Should create a return expression", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 10)
		predicate, err := expressions.NewVariableExpression(sourceMap, "testExp")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)
		So(exp, ShouldHaveSameTypeAs, &expressions.ReturnExpression{})
	})

	Convey("Should not create a return expression with an empty predicate", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 1)
		exp, err := expressions.NewReturnExpression(sourceMap, nil)

		So(err, ShouldBeError)
		So(exp, ShouldBeNil)
	})

	Convey("Should exec a return expression with an existing predicate", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "test")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := rootScope.Fork()
		So(scope.SetVariable("test", values.NewString("value")), ShouldBeNil)
		So(fn(), ShouldBeNil)

		value, err := exp.Exec(context.Background(), scope)
		So(err, ShouldBeNil)
		So(value, ShouldNotBeNil)
		So(value, ShouldEqual, "value")
	})

	Convey("Should not exec a return expression with a missing predicate", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate, err := expressions.NewVariableExpression(sourceMap, "notExist")
		So(err, ShouldBeNil)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		rootScope, fn := core.NewRootScope()
		scope := rootScope.Fork()
		scope.SetVariable("test", values.NewString("value"))
		fn()

		value, err := exp.Exec(context.Background(), scope)
		So(err, ShouldNotBeNil)
		So(err.(*core.SourceErrorDetail).BaseError, ShouldHaveSameTypeAs, core.ErrNotFound)
		So(value, ShouldHaveSameTypeAs, values.None)
	})

	Convey("Should stop an execution when context is cancelled", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 1)
		predicate := literals.NewIntLiteral(1)

		exp, err := expressions.NewReturnExpression(sourceMap, predicate)
		So(err, ShouldBeNil)

		rootScope, _ := core.NewRootScope()
		scope := rootScope.Fork()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err = exp.Exec(ctx, scope)
		So(err, ShouldEqual, core.ErrTerminated)
	})
}
