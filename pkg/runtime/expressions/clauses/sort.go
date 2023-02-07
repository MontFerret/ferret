package clauses

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	SorterExpression struct {
		expression core.Expression
		direction  collections.SortDirection
	}

	SortClause struct {
		src        core.SourceMap
		dataSource collections.Iterable
		sorters    []*SorterExpression
	}
)

func NewSorterExpression(expression core.Expression, direction collections.SortDirection) (*SorterExpression, error) {
	if expression == nil {
		return nil, core.Error(core.ErrMissedArgument, "reducer")
	}

	if !collections.IsValidSortDirection(direction) {
		return nil, core.Error(core.ErrInvalidArgument, "direction")
	}

	return &SorterExpression{expression, direction}, nil
}

func NewSortClause(
	src core.SourceMap,
	dataSource collections.Iterable,
	sorters ...*SorterExpression,
) (collections.Iterable, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "dataSource source")
	}

	if len(sorters) == 0 {
		return nil, core.Error(core.ErrMissedArgument, "sorters")
	}

	return &SortClause{src, dataSource, sorters}, nil
}

func (clause *SortClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	sorters := make([]*collections.Sorter, len(clause.sorters))

	// converting sorter reducer into collections.Sorter
	for idx, srt := range clause.sorters {
		sorter, err := newSorter(srt)

		if err != nil {
			return nil, err
		}

		sorters[idx] = sorter
	}

	return collections.NewSortIterator(src, sorters...)
}

func newSorter(srt *SorterExpression) (*collections.Sorter, error) {
	return collections.NewSorter(func(ctx context.Context, first, second *core.Scope) (int64, error) {
		f, err := srt.expression.Exec(ctx, first)

		if err != nil {
			return -1, err
		}

		s, err := srt.expression.Exec(ctx, second)

		if err != nil {
			return -1, err
		}

		return f.Compare(s), nil
	}, srt.direction)
}
