package expressions_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewParameterExpression(t *testing.T) {
	Convey("Should match", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 1)
		s, err := expressions.NewParameterExpression(sourceMap, "test")

		So(err, ShouldBeNil)
		So(s, ShouldHaveSameTypeAs, &expressions.ParameterExpression{})
	})

	Convey("Should match", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 1)
		s, err := expressions.NewParameterExpression(sourceMap, "")

		So(err, ShouldNotBeNil)
		So(err, ShouldHaveSameTypeAs, core.ErrMissedArgument)
		So(s, ShouldBeNil)
	})
}

func TestParameterExpressionExec(t *testing.T) {
	Convey("Should match", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 10)
		existExp, err := expressions.NewParameterExpression(sourceMap, "param1")

		params := make(map[string]core.Value)
		params["param1"] = values.NewInt(1)

		ctx := core.ParamsWith(context.Background(), params)
		value, err := existExp.Exec(ctx, &core.Scope{})

		So(err, ShouldBeNil)
		So(value.Type(), ShouldEqual, core.IntType)
		So(value.String(), ShouldEqual, "1")
	})

	Convey("Should match", t, func() {
		sourceMap := core.NewSourceMap("test", 1, 10)
		notExistExp, err := expressions.NewParameterExpression(sourceMap, "param2")

		params := make(map[string]core.Value)
		params["param1"] = values.NewInt(1)

		ctx := core.ParamsWith(context.Background(), params)
		value, err := notExistExp.Exec(ctx, &core.Scope{})

		So(err, ShouldNotBeNil)
		So(err, ShouldHaveSameTypeAs, core.ErrNotFound)
		So(value.Type(), ShouldEqual, core.NoneType)
	})
}
