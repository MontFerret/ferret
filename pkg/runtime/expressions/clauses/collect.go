package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Collect struct {
		Grouping *CollectGrouping
	}

	CollectGrouping struct {
		Selectors  []*CollectSelector
		Projection *CollectProjection
		Count      *CollectCount
	}

	CollectCount struct {
		variable string
	}

	CollectProjection struct {
		selector *CollectSelector
	}

	CollectClause struct {
		src        core.SourceMap
		dataSource collections.Iterable
		params     *Collect
	}
)

func NewCollectCount(variable string) *CollectCount {
	return &CollectCount{variable}
}

func NewCollectProjection(selector *CollectSelector) *CollectProjection {
	return &CollectProjection{selector}
}

func NewCollectClause(
	src core.SourceMap,
	dataSource collections.Iterable,
	params *Collect,
) (collections.Iterable, error) {
	if dataSource == nil {
		return nil, core.Error(core.ErrMissedArgument, "dataSource source")
	}

	return &CollectClause{src, dataSource, params}, nil
}

func (clause *CollectClause) Variables() collections.Variables {
	vars := make(collections.Variables, 0, 10)
	grouping := clause.params.Grouping

	if grouping != nil {
		if grouping.Selectors != nil {
			for _, selector := range grouping.Selectors {
				vars = append(vars, selector.variable)
			}
		}

		if grouping.Projection != nil {
			vars = append(vars, clause.params.Grouping.Projection.selector.variable)
		}

		if grouping.Count != nil {
			vars = append(vars, clause.params.Grouping.Count.variable)
		}
	}

	return vars
}

func (clause *CollectClause) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	srcIterator, err := clause.dataSource.Iterate(ctx, scope)

	if err != nil {
		return nil, core.SourceError(clause.src, err)
	}

	srcVariables := clause.dataSource.Variables()

	return NewCollectIterator(
		clause.src,
		clause.params,
		srcIterator,
		srcVariables,
		ctx,
		scope,
	)
}
