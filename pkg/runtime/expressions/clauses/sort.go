package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/source"
)

type (
	SorterExpression struct {
		expression core.Expression
		direction  collections.SortDirection
	}
	SortClause struct {
		src          core.SourceMap
		dataSource   source.DataSource
		variableName string
		sorters      []*SorterExpression
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
	dataSource source.DataSource,
	variableName string,
	sorters ...*SorterExpression,
) *SortClause {
	return &SortClause{src, dataSource, variableName, sorters}
}

func (clause *SortClause) Variables() []string {
	return clause.dataSource.Variables()
}

func (clause *SortClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	sorters := make([]*collections.Sorter, len(clause.sorters))

	// converting sorter expression into collections.Sorter
	for idx, srt := range clause.sorters {
		sorter, err := collections.NewSorter(func(first core.Value, second core.Value) (int, error) {
			scope1 := scope.Fork()
			scope1.SetVariable(clause.variableName, first)

			f, err := srt.expression.Exec(ctx, scope1)

			if err != nil {
				return -1, err
			}

			scope2 := scope.Fork()
			scope2.SetVariable(clause.variableName, second)

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
