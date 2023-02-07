package clauses

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type FilterClause struct {
	src        core.SourceMap
	dataSource collections.Iterable
	predicate  core.Expression
}

func NewFilterClause(
	src core.SourceMap,
	dataSource collections.Iterable,
	predicate core.Expression,
) (collections.Iterable, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "dataSource source")
	}

	if predicate == nil {
		return nil, core.Error(core.ErrMissedArgument, "predicate")
	}

	return &FilterClause{
		src, dataSource,
		predicate,
	}, nil
}

func (clause *FilterClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewFilterIterator(src, clause.filter)
}

func (clause *FilterClause) filter(ctx context.Context, scope *core.Scope) (bool, error) {
	ret, err := clause.predicate.Exec(ctx, scope)

	if err != nil {
		return false, err
	}

	if ret == values.True {
		return true, nil
	}

	return false, nil
}
