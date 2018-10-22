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
	dataSource collections.Iterable
	predicate  core.Expression
	distinct   bool
	spread     bool
}

func NewForExpression(
	src core.SourceMap,
	dataSource collections.Iterable,
	predicate core.Expression,
	distinct,
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
		distinct,
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

func (e *ForExpression) AddCollect(src core.SourceMap, params *clauses.CollectParams) error {
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

	// Hash map for a check for uniqueness
	var hashTable map[uint64]bool

	if e.distinct {
		hashTable = make(map[uint64]bool)
	}

	res := values.NewArray(10)
	variables := e.dataSource.Variables()

	for iterator.HasNext() {
		ds, err := iterator.Next()

		if err != nil {
			return values.None, core.SourceError(e.src, err)
		}

		innerScope := scope.Fork()

		if err := ds.Apply(innerScope, variables); err != nil {
			return values.None, err
		}

		out, err := e.predicate.Exec(ctx, innerScope)

		if err != nil {
			return values.None, err
		}

		var add bool

		// The result shouldn't be distinct
		// Just add the output
		if !e.distinct {
			add = true
		} else {
			// We need to check whether the value already exists in the result set
			hash := out.Hash()
			_, exists := hashTable[hash]

			if !exists {
				hashTable[hash] = true
				add = true
			}
		}

		if add {
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
	}

	return res, nil
}
