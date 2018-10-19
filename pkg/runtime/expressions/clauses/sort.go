package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/datasource"
)

type (
	SorterExpression struct {
		expression core.Expression
		direction  collections.SortDirection
	}
	SortClause struct {
		src        core.SourceMap
		dataSource datasource.DataSource
		sorters    []*SorterExpression
	}
)

func NewSorterExpression(expression core.Expression, direction collections.SortDirection) (*SorterExpression, error) {
	if expression == nil {
		return nil, core.Error(core.ErrMissedArgument, "expression")
	}

	if !collections.IsValidSortDirection(direction) {
		return nil, core.Error(core.ErrInvalidArgument, "direction")
	}

	return &SorterExpression{expression, direction}, nil
}

func NewSortClause(
	src core.SourceMap,
	dataSource datasource.DataSource,
	sorters ...*SorterExpression,
) *SortClause {
	return &SortClause{src, dataSource, sorters}
}

func (clause *SortClause) Variables() datasource.Variables {
	return clause.dataSource.Variables()
}

func (clause *SortClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	sorters := make([]*collections.Sorter, len(clause.sorters))
	variables := clause.dataSource.Variables()

	// converting sorter expression into collections.Sorter
	for idx, srt := range clause.sorters {
		sorter, err := collections.NewSorter(func(first collections.ResultSet, second collections.ResultSet) (int, error) {
			scope1 := scope.Fork()

			if err := variables.Apply(scope1, first); err != nil {
				return -1, err
			}

			f, err := srt.expression.Exec(ctx, scope1)

			if err != nil {
				return -1, err
			}

			scope2 := scope.Fork()

			if err := variables.Apply(scope2, second); err != nil {
				return -1, err
			}

			s, err := srt.expression.Exec(ctx, scope2)

			if err != nil {
				return -1, err
			}

			return f.Compare(s), nil
		}, srt.direction)

		if err != nil {
			return nil, err
		}

		sorters[idx] = sorter
	}

	return collections.NewSortIterator(src, sorters...)
}
