package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ForExpression struct {
	src         core.SourceMap
	dataSource  collections.Iterable
	predicate   core.Expression
	distinct    bool
	spread      bool
	passThrough bool
}

func NewForExpression(
	src core.SourceMap,
	dataSource collections.Iterable,
	predicate core.Expression,
	distinct,
	spread,
	passThrough bool,
) (*ForExpression, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "missed source expression")
	}

	if predicate == nil {
		return nil, core.Error(core.ErrMissedArgument, "missed return expression")
	}

	exp := new(ForExpression)
	exp.src = src
	exp.dataSource = dataSource
	exp.predicate = predicate
	exp.distinct = distinct
	exp.spread = spread
	exp.passThrough = passThrough

	return exp, nil
}

func (e *ForExpression) AddLimit(src core.SourceMap, size, count core.Expression) error {
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

func (e *ForExpression) AddCollect(src core.SourceMap, params *clauses.Collect) error {
	collect, err := clauses.NewCollectClause(src, e.dataSource, params)

	if err != nil {
		return err
	}

	e.dataSource = collect

	return nil
}

func (e *ForExpression) AddStatement(stmt core.Expression) error {
	tap, ok := e.dataSource.(*BlockExpression)

	if !ok {
		t, err := NewBlockExpression(e.dataSource)

		if err != nil {
			return err
		}

		tap = t
		e.dataSource = tap
	}

	tap.Add(stmt)

	return nil
}

func (e *ForExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	select {
	case <-ctx.Done():
		return values.None, core.ErrTerminated
	default:
		iterator, err := e.dataSource.Iterate(ctx, scope)

		if err != nil {
			return values.None, err
		}

		res := NewForResult(10).
			Distinct(e.distinct).
			Spread(e.spread).
			PassThrough(e.passThrough)

		for {
			nextScope, err := iterator.Next(ctx, scope)

			if err != nil {
				if core.IsNoMoreData(err) {
					break
				}

				return values.None, core.SourceError(e.src, err)
			}

			out, err := e.predicate.Exec(ctx, nextScope)

			if err != nil {
				return values.None, err
			}

			res.Push(out)
		}

		return res.ToArray(), nil
	}
}
