package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Collect struct {
		Group     *CollectGroup
		Count     *CollectCount
		Aggregate *CollectAggregate
	}

	CollectGroup struct {
		Selectors  []*CollectSelector
		Projection *CollectProjection
		Count      *CollectCount
		Aggregate  *CollectAggregate
	}

	CollectCount struct {
		variable string
	}

	CollectProjection struct {
		selector *CollectSelector
	}

	CollectAggregate struct {
		Selectors []*CollectAggregateSelector
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

	if clause.params.Group != nil {
		grouping := clause.params.Group

		for _, selector := range grouping.Selectors {
			vars = append(vars, selector.variable)
		}

		if grouping.Projection != nil {
			vars = append(vars, clause.params.Group.Projection.selector.variable)
		}

		if grouping.Count != nil {
			vars = append(vars, clause.params.Group.Count.variable)
		}

		if grouping.Aggregate != nil {
			for _, selector := range grouping.Aggregate.Selectors {
				vars = append(vars, selector.variable)
			}
		}
	} else if clause.params.Count != nil {
		vars = append(vars, clause.params.Count.variable)
	} else if clause.params.Aggregate != nil {
		for _, selector := range clause.params.Aggregate.Selectors {
			vars = append(vars, selector.variable)
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
