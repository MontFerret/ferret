package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/datasource"
)

type DistinctClause struct {
	src        core.SourceMap
	dataSource datasource.DataSource
}

func NewDistinctClause(
	src core.SourceMap,
	dataSource datasource.DataSource,
) *DistinctClause {
	return &DistinctClause{src, dataSource}
}

func (clause *DistinctClause) Variables() datasource.Variables {
	return clause.dataSource.Variables()
}

func (clause *DistinctClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewUniqueIterator(src)
}
