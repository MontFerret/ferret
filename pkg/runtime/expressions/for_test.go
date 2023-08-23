package expressions

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/literals"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/operators"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

func testIntArrayLiteral() (collections.Iterable, string) {
	dataSource, _ := NewForInIterableExpression(
		core.SourceMap{},
		"val",
		"",
		literals.NewArrayLiteralWith([]core.Expression{
			literals.NewIntLiteral(0),
			literals.NewIntLiteral(1),
			literals.NewIntLiteral(2),
			literals.NewIntLiteral(3),
			literals.NewIntLiteral(4),
		}),
	)

	return dataSource, "val"
}

func TestNewForExpression(t *testing.T) {
	dataSource, _ := testIntArrayLiteral()
	returnExp, _ := NewVariableExpression(core.SourceMap{}, "testExp")

	Convey("NewForExpression", t, func() {
		Convey("should return new ForExpresssion.", func() {
			forExp, err := NewForExpression(core.SourceMap{}, dataSource, returnExp, false, false, false)
			So(forExp, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		Convey("should return error when a dataSource is nil", func() {
			forExp, err := NewForExpression(core.SourceMap{}, nil, returnExp, false, false, false)
			So(forExp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, err)
		})

		Convey("should return error when a predicate is nil", func() {
			forExp, err := NewForExpression(core.SourceMap{}, dataSource, nil, false, false, false)
			So(forExp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, err)
		})
	})
}

func TestAddLimit(t *testing.T) {
	dataSource, valName := testIntArrayLiteral()

	returnInValExp, _ := NewVariableExpression(core.SourceMap{}, valName)
	returnExp, _ := NewReturnExpression(
		core.SourceMap{},
		returnInValExp,
	)

	Convey("AddLimit", t, func() {
		Convey("should success. (An Error should be nil.)", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)
			err := forExp.AddLimit(core.SourceMap{}, literals.NewIntLiteral(3), literals.NewIntLiteral(0))
			So(err, ShouldBeNil)
		})

		Convey("should return an error.", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)
			forExp.dataSource = nil
			err := forExp.AddLimit(core.SourceMap{}, literals.NewIntLiteral(3), literals.NewIntLiteral(0))
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddFilter(t *testing.T) {
	dataSource, valName := testIntArrayLiteral()

	returnInValExp, _ := NewVariableExpression(core.SourceMap{}, valName)
	returnExp, _ := NewReturnExpression(
		core.SourceMap{},
		returnInValExp,
	)

	testFilter := func() core.Expression {
		filterLeftExp, _ := NewVariableExpression(core.SourceMap{}, valName)
		filterRightExp := literals.NewIntLiteral(3)
		filterExp, _ := operators.NewEqualityOperator(core.SourceMap{}, filterLeftExp, filterRightExp, "<")
		return filterExp
	}

	Convey("AddFilter", t, func() {
		Convey("should success. (An Error should be nil.)", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)

			err := forExp.AddFilter(core.SourceMap{}, testFilter())
			So(err, ShouldBeNil)
		})

		Convey("should return an error.", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)
			forExp.dataSource = nil
			err := forExp.AddFilter(core.SourceMap{}, testFilter())
			So(err, ShouldNotBeNil)
		})
	})

}

func TestAddSort(t *testing.T) {
	dataSource, valName := testIntArrayLiteral()

	returnInValExp, _ := NewVariableExpression(core.SourceMap{}, valName)
	returnExp, _ := NewReturnExpression(
		core.SourceMap{},
		returnInValExp,
	)

	testSort := func() *clauses.SorterExpression {
		valExp, _ := NewVariableExpression(core.SourceMap{}, valName)
		sortExp, _ := clauses.NewSorterExpression(valExp, collections.SortDirectionDesc)
		return sortExp
	}

	Convey("AddSort", t, func() {
		Convey("should success.(An Error should be nil.)", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)

			err := forExp.AddSort(core.SourceMap{}, testSort())
			So(err, ShouldBeNil)
		})

		Convey("should return an error.", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)
			forExp.dataSource = nil
			err := forExp.AddSort(core.SourceMap{}, testSort())
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddCollect(t *testing.T) {
	dataSource, valName := testIntArrayLiteral()

	selectorValName := "selector"

	testCollect := func() *clauses.Collect {
		eleVal, _ := NewVariableExpression(core.SourceMap{}, valName)
		selector, _ := clauses.NewCollectSelector(selectorValName, eleVal)
		collect, _ := clauses.NewCollect(
			[]*clauses.CollectSelector{selector}, nil, nil, nil,
		)

		return collect
	}

	returnInValExp, _ := NewVariableExpression(core.SourceMap{}, selectorValName)
	returnExp, _ := NewReturnExpression(
		core.SourceMap{},
		returnInValExp,
	)

	Convey("AddCollect", t, func() {
		Convey("should success. (Error should be nil.)", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)

			err := forExp.AddCollect(core.SourceMap{}, testCollect())
			So(err, ShouldBeNil)
		})

		Convey("should return an error.", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)
			forExp.dataSource = nil
			err := forExp.AddCollect(core.SourceMap{}, testCollect())
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddStatement(t *testing.T) {
	dataSource, valName := testIntArrayLiteral()

	returnInValExp, _ := NewVariableExpression(core.SourceMap{}, valName)
	returnExp, _ := NewReturnExpression(
		core.SourceMap{},
		returnInValExp,
	)

	testStatement, _ := NewVariableDeclarationExpression(core.SourceMap{}, "newVal", literals.NewIntLiteral(0))

	Convey("AddStatement", t, func() {
		Convey("should success.", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)

			err := forExp.AddStatement(testStatement)
			So(err, ShouldBeNil)
		})

		Convey("should return an error. (Error should be nil.)", func() {
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)
			forExp.dataSource = nil
			err := forExp.AddStatement(testStatement)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestExec(t *testing.T) {
	Convey("Exec", t, func() {
		Convey("should success.", func() {
			dataSource, valName := testIntArrayLiteral()

			returnInValExp, _ := NewVariableExpression(core.SourceMap{}, valName)
			returnExp, _ := NewReturnExpression(
				core.SourceMap{},
				returnInValExp,
			)

			rootScope, closeFn := core.NewRootScope()
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)

			resultVal, err := forExp.Exec(context.Background(), rootScope)
			So(err, ShouldBeNil)
			resultArr, ok := resultVal.(*values.Array)
			So(ok, ShouldBeTrue)
			So(resultArr.Length(), ShouldEqual, values.NewInt(5))
			compareArr := values.NewArrayOf([]core.Value{
				values.NewInt(0),
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
			})
			So(resultArr.Compare(compareArr), ShouldEqual, 0)

			closeFn()
		})

		Convey("should stop an execution when context is cancelled.", func() {
			dataSource, valName := testIntArrayLiteral()

			returnInValExp, _ := NewVariableExpression(core.SourceMap{}, valName)
			returnExp, _ := NewReturnExpression(
				core.SourceMap{},
				returnInValExp,
			)

			rootScope, closeFn := core.NewRootScope()
			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)

			ctx0 := context.Background()
			ctx1, cancelFn := context.WithCancel(ctx0)
			cancelFn()

			_, err := forExp.Exec(ctx1, rootScope)
			So(err, ShouldNotBeNil)

			closeFn()
		})

		Convey("should return an error when a dataSource expression is invalidated.", func() {
			rootScope, closeFn := core.NewRootScope()
			errorDataSource, valName := func() (collections.Iterable, string) {
				dataSource, _ := NewForInIterableExpression(
					core.SourceMap{},
					"val",
					"",
					literals.NewIntLiteral(1),
				)

				return dataSource, "val"
			}()

			returnInValExp, _ := NewVariableExpression(core.SourceMap{}, valName)
			returnExp, _ := NewReturnExpression(
				core.SourceMap{},
				returnInValExp,
			)

			forExp, _ := NewForExpression(
				core.SourceMap{},
				errorDataSource,
				returnExp,
				false,
				false,
				false,
			)

			result, err := forExp.Exec(context.Background(), rootScope)
			So(result, ShouldEqual, values.None)
			So(err, ShouldNotBeNil)

			closeFn()
		})

		Convey("should return an error when element expressions of dataSource is invalidated.", func() {
			rootScope, closeFn := core.NewRootScope()

			errorDataSource, valName := func() (collections.Iterable, string) {
				errorEle, _ := NewVariableExpression(core.SourceMap{}, "a")
				dataSource, _ := NewForInIterableExpression(
					core.SourceMap{},
					"val",
					"",
					literals.NewArrayLiteralWith([]core.Expression{
						literals.NewIntLiteral(0),
						literals.NewIntLiteral(1),
						literals.NewIntLiteral(2),
						errorEle,
					}),
				)

				return dataSource, "val"
			}()

			returnInValExp, _ := NewVariableExpression(core.SourceMap{}, valName)
			returnExp, _ := NewReturnExpression(
				core.SourceMap{},
				returnInValExp,
			)

			forExp, _ := NewForExpression(
				core.SourceMap{},
				errorDataSource,
				returnExp,
				false,
				false,
				false,
			)

			result, err := forExp.Exec(context.Background(), rootScope)
			So(result, ShouldEqual, values.None)
			So(err, ShouldNotBeNil)

			closeFn()
		})

		Convey("should return an error when an predicate expression is invalidated.", func() {
			rootScope, closeFn := core.NewRootScope()
			dataSource, _ := testIntArrayLiteral()

			returnInValExp, _ := NewVariableExpression(core.SourceMap{}, "notExistVal")
			returnExp, _ := NewReturnExpression(
				core.SourceMap{},
				returnInValExp,
			)

			forExp, _ := NewForExpression(
				core.SourceMap{},
				dataSource,
				returnExp,
				false,
				false,
				false,
			)

			result, err := forExp.Exec(context.Background(), rootScope)
			So(result, ShouldEqual, values.None)
			So(err, ShouldNotBeNil)

			closeFn()
		})
	})
}
