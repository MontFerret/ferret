package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type LimitClause struct {
	src        core.SourceMap
	dataSource collections.Iterable
	count      int
	offset     int
}

func NewLimitClause(
	src core.SourceMap,
	dataSource collections.Iterable,
	count int,
	offset int,
) (collections.Iterable, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "dataSource source")
	}

	return &LimitClause{src, dataSource, count, offset}, nil
}

func (clause *LimitClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, core.SourceError(clause.src, err)
	}

	iterator, err := collections.NewLimitIterator(
		src,
		clause.count,
		clause.offset,
	)

	if err != nil {
		return nil, core.SourceError(clause.src, err)
	}

	return iterator, nil
}
