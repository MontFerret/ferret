package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/source"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

type ForExpression struct {
	src        core.SourceMap
	valVar     string
	keyVar     string
	dataSource source.DataSource
	predicate  core.Expression
	spread     bool
}

func NewForExpression(
	src core.SourceMap,
	valVar string,
	keyVar string,
	dataSource source.DataSource,
	predicate core.Expression,
	spread bool,
) (*ForExpression, error) {
	if valVar == "" {
		return nil, errors.Wrap(core.ErrInvalidArgument, "valVar is empty")
	}

	if core.IsNil(dataSource) {
		return nil, errors.Wrap(core.ErrMissedArgument, "missed source expression")
	}

	if core.IsNil(predicate) {
		return nil, errors.Wrap(core.ErrMissedArgument, "missed return expression")
	}

	return &ForExpression{
		src,
		valVar, keyVar,
		dataSource,
		predicate,
		spread,
	}, nil
}

func (e *ForExpression) AddLimit(src core.SourceMap, size, count int) {
	e.dataSource = clauses.NewLimitClause(src, e.dataSource, size, count)
}

func (e *ForExpression) AddFilter(src core.SourceMap, exp core.Expression) {
	e.dataSource = clauses.NewFilterClause(src, e.dataSource, e.valVar, e.keyVar, exp)
}

func (e *ForExpression) AddSort(src core.SourceMap, sorters ...*clauses.SorterExpression) {
	e.dataSource = clauses.NewSortClause(src, e.dataSource, e.valVar, sorters...)
}

func (e *ForExpression) AddDistinct(src core.SourceMap) {
	e.dataSource = clauses.NewDistinctClause(src, e.dataSource)
}

func (e *ForExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	iterator, err := e.dataSource.Iterate(ctx, scope)

	if err != nil {
		return values.None, err
	}

	res := values.NewArray(10)

	for iterator.HasNext() {
		val, key, err := iterator.Next()

		if err != nil {
			return values.None, core.SourceError(e.src, err)
		}

		innerScope := scope.Fork()
		innerScope.SetVariable(e.valVar, val)

		if e.keyVar != "" {
			innerScope.SetVariable(e.keyVar, key)
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
