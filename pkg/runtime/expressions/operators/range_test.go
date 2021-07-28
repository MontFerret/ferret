package operators_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/operators"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

type MockExpressionInt int

func (e MockExpressionInt) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return values.NewInt(int(e)), nil
}

type MockExpressionFloat float64

func (e MockExpressionFloat) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return values.NewFloat(float64(e)), nil
}

type MockExpressionString string

func (e MockExpressionString) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return values.NewString(string(e)), nil
}

type FailingMockExpression int

func (e FailingMockExpression) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return nil, errors.New("MockError")
}

var sourceMap = core.NewSourceMap("hello", 2, 3)
var rootScope, _ = core.NewRootScope()

func TestRangeOperator(t *testing.T) {
	var correctEx MockExpressionInt = 10
	var failEx FailingMockExpression = 0
	Convey("NewRangeOperator", t, func() {
		Convey("Should return *RangeOperator, given correct arguments", func() {
			res, _ := operators.NewRangeOperator(sourceMap, correctEx, correctEx)
			So(res, ShouldHaveSameTypeAs, &operators.RangeOperator{})
		})
		Convey("Should return error if no left operator is given", func() {
			_, err := operators.NewRangeOperator(sourceMap, nil, correctEx)
			So(err, ShouldHaveSameTypeAs, core.ErrMissedArgument)
		})
		Convey("Should return error if no right operator is given", func() {
			_, err := operators.NewRangeOperator(sourceMap, correctEx, nil)
			So(err, ShouldHaveSameTypeAs, core.ErrMissedArgument)
		})
	})
	Convey(".Exec", t, func() {
		Convey("Should return error if left expression fails", func() {
			rangeOp, _ := operators.NewRangeOperator(sourceMap, failEx, correctEx)
			_, err := rangeOp.Exec(context.TODO(), rootScope)
			So(err, ShouldHaveSameTypeAs, &core.SourceErrorDetail{})
		})
		Convey("Should return error if right expression fails", func() {
			rangeOp, _ := operators.NewRangeOperator(sourceMap, correctEx, failEx)
			_, err := rangeOp.Exec(context.TODO(), rootScope)
			So(err, ShouldHaveSameTypeAs, &core.SourceErrorDetail{})
		})
	})
	Convey(".Eval", t, func() {
		var stringEx MockExpressionString = "noInt/noFloat"
		Convey("Should return error if given non int/float left Expression", func() {
			rangeOp, _ := operators.NewRangeOperator(sourceMap, stringEx, correctEx)
			_, err := rangeOp.Exec(context.TODO(), rootScope)
			So(err, ShouldHaveSameTypeAs, &core.SourceErrorDetail{})
		})
		Convey("Should return error if given non int/float right Expression", func() {
			rangeOp, _ := operators.NewRangeOperator(sourceMap, correctEx, stringEx)
			_, err := rangeOp.Exec(context.TODO(), rootScope)
			So(err, ShouldHaveSameTypeAs, &core.SourceErrorDetail{})
		})
		Convey("Should return Array with range [start,...end] given ints", func() {
			var start MockExpressionInt = 1
			var end MockExpressionInt = 14

			rangeOp, _ := operators.NewRangeOperator(sourceMap, start, end)
			arr, err := rangeOp.Exec(context.TODO(), rootScope)
			So(err, ShouldHaveSameTypeAs, nil)
			So(arr, ShouldHaveSameTypeAs, &values.Array{})
			So(arr.String(), ShouldEqual, "[1,2,3,4,5,6,7,8,9,10,11,12,13,14]")
		})
		Convey("Should return Array with range [start,...end] given floats", func() {
			var start MockExpressionFloat = -2.2
			var end MockExpressionFloat = 12.4

			rangeOp, _ := operators.NewRangeOperator(sourceMap, start, end)
			arr, err := rangeOp.Exec(context.TODO(), rootScope)
			So(err, ShouldHaveSameTypeAs, nil)
			So(arr, ShouldHaveSameTypeAs, &values.Array{})
			So(arr.String(), ShouldEqual, "[-2,-1,0,1,2,3,4,5,6,7,8,9,10,11,12]")
		})
	})
}
