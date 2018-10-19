package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/datasource"
)

type LimitClause struct {
	src        core.SourceMap
	dataSource datasource.DataSource
	count      int
	offset     int
}

func NewLimitClause(
	src core.SourceMap,
	dataSource datasource.DataSource,
	count int,
	offset int,
) (datasource.DataSource, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "dataSource source")
	}

	return &LimitClause{src, dataSource, count, offset}, nil
}

func (clause *LimitClause) Variables() datasource.Variables {
	return clause.dataSource.Variables()
}

func (clause *LimitClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	src, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewLimitIterator(src, clause.count, clause.offset)
}
