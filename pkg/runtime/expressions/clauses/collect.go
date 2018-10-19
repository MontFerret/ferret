package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/datasource"
)

type (
	CollectParams struct {
		Grouping []*CollectSelector
	}

	CollectClause struct {
		src        core.SourceMap
		dataSource datasource.DataSource
		params     CollectParams
	}
)

func NewCollect(
	src core.SourceMap,
	dataSource datasource.DataSource,
	params CollectParams,
) (datasource.DataSource, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "dataSource source")
	}

	return &CollectClause{src, dataSource, params}, nil
}

func (clause *CollectClause) Variables() datasource.Variables {
	vars := make(datasource.Variables, 0, len(clause.params.Grouping))

	for _, selector := range clause.params.Grouping {
		vars = append(vars, selector.variable)
	}

	return vars
}

func (clause *CollectClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	srcIterator, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, core.SourceError(clause.src, err)
	}

	srcVariables := clause.dataSource.Variables()

	return NewCollectGroupingIterator(
		clause.src,
		clause.params.Grouping,
		srcIterator,
		srcVariables,
		ctx,
		scope,
	), nil
}
