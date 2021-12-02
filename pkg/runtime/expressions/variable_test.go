package expressions_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var sourceMap = core.NewSourceMap("hello", 2, 3)
var rootScope, _ = core.NewRootScope()
var _ = rootScope.SetVariable("key", values.NewString("value"))

func TestNewVariableExpression(t *testing.T) {

	Convey("Should not throw error and create a VariableExpression given correct arguments", t, func() {
		ret, err := expressions.NewVariableExpression(sourceMap, "foo")

		So(ret, ShouldHaveSameTypeAs, &expressions.VariableExpression{})
		So(err, ShouldBeNil)
	})

	Convey("Should throw error when given no name argument", t, func() {
		_, err := expressions.NewVariableExpression(sourceMap, "")

		So(err, ShouldHaveSameTypeAs, core.ErrMissedArgument)
	})

	Convey("Calling .Exec should return the correct variable set in the given scope", t, func() {
		ret, _ := expressions.NewVariableExpression(sourceMap, "key")
		value, err := ret.Exec(context.TODO(), rootScope)

		So(value, ShouldEqual, "value")
		So(err, ShouldBeNil)
	})
}

func TestNewVariableDeclarationExpression(t *testing.T) {
	Convey("Should not throw error and create a NewVariableDeclarationExpression given correct arguments", t, func() {
		variableExpression, _ := expressions.NewVariableExpression(sourceMap, "foo")
		ret, err := expressions.NewVariableDeclarationExpression(sourceMap, "test", variableExpression)

		So(ret, ShouldHaveSameTypeAs, &expressions.VariableDeclarationExpression{})
		So(err, ShouldBeNil)
	})

	Convey("Should throw error if init argument is nil", t, func() {
		_, err := expressions.NewVariableDeclarationExpression(sourceMap, "test", nil)

		So(err, ShouldHaveSameTypeAs, core.ErrMissedArgument)
	})

	Convey("Calling .Exec should add the value retrieved by its VariableExpression with its own name as key to the given scope", t, func() {
		variableExpression, _ := expressions.NewVariableExpression(sourceMap, "key")
		variableDeclarationExpression, _ := expressions.NewVariableDeclarationExpression(sourceMap, "keyTwo", variableExpression)
		_, err := variableDeclarationExpression.Exec(context.TODO(), rootScope)

		So(err, ShouldBeNil)

		value, _ := rootScope.GetVariable("key")
		value2, _ := rootScope.GetVariable("keyTwo")

		So(value, ShouldEqual, "value")
		So(value2, ShouldEqual, "value")
	})
}
