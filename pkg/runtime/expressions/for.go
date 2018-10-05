package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/clauses"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

type ForExpression struct {
	src        core.SourceMap
	valVar     string
	keyVar     string
	dataSource collections.IterableExpression
	predicate  core.Expression
	distinct   bool
	spread     bool
}

func NewForExpression(
	src core.SourceMap,
	valVar string,
	keyVar string,
	dataSource collections.IterableExpression,
	predicate core.Expression,
	distinct bool,
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
		distinct,
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

func (e *ForExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	iterator, err := e.dataSource.Iterate(ctx, scope)

	if err != nil {
		return values.None, err
	}

	// Hash map for a check for uniqueness
	var hashes map[uint64]bool

	if e.distinct {
		hashes = make(map[uint64]bool)
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

		var el core.Value

		// The result shouldn't be distinct
		// Just add the output
		if !e.distinct {
			el = out
		} else {
			// We need to check whether the value already exists in the result set
			hash := out.Hash()
			_, exists := hashes[hash]

			if !exists {
				hashes[hash] = true
				el = out
			}
		}

		if el != nil {
			if !e.spread {
				res.Push(el)
			} else {
				elements := el.(*values.Array)

				elements.ForEach(func(i core.Value, _ int) bool {
					res.Push(i)

					return true
				})
			}
		}
	}

	return res, nil
}
