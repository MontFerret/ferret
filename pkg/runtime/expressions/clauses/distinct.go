package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type DistinctClause struct {
	src        core.SourceMap
	dataSource collections.DataSource
}

func NewDistinctClause(
	src core.SourceMap,
	dataSource collections.DataSource,
) *DistinctClause {
	return &DistinctClause{src, dataSource}
}

func (clause *DistinctClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewUniqueIterator(src)
}
