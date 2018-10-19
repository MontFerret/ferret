package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/datasource"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

type ForExpression struct {
	src        core.SourceMap
	dataSource datasource.DataSource
	predicate  core.Expression
	spread     bool
}

func NewForExpression(
	src core.SourceMap,
	dataSource datasource.DataSource,
	predicate core.Expression,
	spread bool,
) (*ForExpression, error) {
	if core.IsNil(dataSource) {
		return nil, errors.Wrap(core.ErrMissedArgument, "missed source expression")
	}

	if core.IsNil(predicate) {
		return nil, errors.Wrap(core.ErrMissedArgument, "missed return expression")
	}

	return &ForExpression{
		src,
		dataSource,
		predicate,
		spread,
	}, nil
}

func (e *ForExpression) AddLimit(src core.SourceMap, size, count int) error {
	limit, err := clauses.NewLimitClause(src, e.dataSource, size, count)

	if err != nil {
		return err
	}

	e.dataSource = limit

	return nil
}

func (e *ForExpression) AddFilter(src core.SourceMap, exp core.Expression) error {
	filter, err := clauses.NewFilterClause(src, e.dataSource, exp)

	if err != nil {
		return err
	}

	e.dataSource = filter

	return nil
}

func (e *ForExpression) AddSort(src core.SourceMap, sorters ...*clauses.SorterExpression) error {
	sort, err := clauses.NewSortClause(src, e.dataSource, sorters...)

	if err != nil {
		return err
	}

	e.dataSource = sort

	return nil
}

func (e *ForExpression) AddDistinct(src core.SourceMap) error {
	distinct, err := clauses.NewDistinctClause(src, e.dataSource)

	if err != nil {
		return err
	}

	e.dataSource = distinct

	return nil
}

func (e *ForExpression) AddCollect(src core.SourceMap, params clauses.CollectParams) error {
	collect, err := clauses.NewCollect(src, e.dataSource, params)

	if err != nil {
		return err
	}

	e.dataSource = collect

	return nil
}

func (e *ForExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	iterator, err := e.dataSource.Iterate(ctx, scope)

	if err != nil {
		return values.None, err
	}

	res := values.NewArray(10)
	variables := e.dataSource.Variables()

	for iterator.HasNext() {
		set, err := iterator.Next()

		if err != nil {
			return values.None, core.SourceError(e.src, err)
		}

		innerScope := scope.Fork()

		// assign returned values to variables for the nested scope
		if err := variables.Apply(innerScope, set); err != nil {
			return values.None, core.SourceError(e.src, err)
		}

		out, err := e.predicate.Exec(ctx, innerScope)

		if err != nil {
			return values.None, err
		}

		if !e.spread {
			res.Push(out)
		} else {
			elements, ok := out.(*values.Array)

			if !ok {
				return values.None, core.Error(core.ErrInvalidOperation, "spread of non-array value")
			}

			elements.ForEach(func(i core.Value, _ int) bool {
				res.Push(i)

				return true
			})
		}
	}

	return res, nil
}
