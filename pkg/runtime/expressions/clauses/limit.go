package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/source"
)

type LimitClause struct {
	src        core.SourceMap
	dataSource source.DataSource
	count      int
	offset     int
}

func NewLimitClause(
	src core.SourceMap,
	dataSource source.DataSource,
	count int,
	offset int,
) *LimitClause {
	return &LimitClause{src, dataSource, count, offset}
}

func (clause *LimitClause) Variables() source.Variables {
	return clause.dataSource.Variables()
}

func (clause *LimitClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewLimitIterator(src, clause.count, clause.offset)
}
