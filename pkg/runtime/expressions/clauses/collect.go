package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Collect struct {
		group     *CollectGroup
		count     *CollectCount
		aggregate *CollectAggregate
	}

	CollectGroup struct {
		selectors  []*CollectSelector
		projection *CollectProjection
		count      *CollectCount
		aggregate  *CollectAggregate
	}

	CollectCount struct {
		variable string
	}

	CollectProjection struct {
		selector *CollectSelector
	}

	CollectAggregate struct {
		selectors []*CollectAggregateSelector
	}

	CollectClause struct {
		src        core.SourceMap
		dataSource collections.Iterable
		params     *Collect
	}
)

func NewCollect(
	selectors []*CollectSelector,
	projection *CollectProjection,
	count *CollectCount,
	aggregate *CollectAggregate,
) (*Collect, error) {
	collect := new(Collect)

	// grouping
	if selectors != nil {
		collect.group = new(CollectGroup)
		collect.group.selectors = selectors

		if projection == nil && count == nil && aggregate == nil {
			return collect, nil
		}

		if projection != nil && count == nil && aggregate == nil {
			collect.group.projection = projection
		} else if projection == nil && count != nil && aggregate == nil {
			collect.group.count = count
		} else if projection == nil && count == nil && aggregate != nil {
			collect.group.aggregate = aggregate
		} else {
			return nil, core.Error(core.ErrInvalidOperation, "projection, count and aggregate cannot be used together")
		}

		return collect, nil
	}

	if count == nil && aggregate != nil {
		collect.aggregate = aggregate
	} else if count != nil && aggregate == nil {
		collect.count = count
	} else {
		return nil, core.Error(core.ErrInvalidOperation, "count and aggregate cannot be used together")
	}

	return collect, nil
}

func NewCollectCount(variable string) (*CollectCount, error) {
	if variable == "" {
		return nil, core.Error(core.ErrMissedArgument, "count variable")
	}

	return &CollectCount{variable}, nil
}

func NewCollectProjection(selector *CollectSelector) (*CollectProjection, error) {
	if selector == nil {
		return nil, core.Error(core.ErrMissedArgument, "projection selector")
	}

	return &CollectProjection{selector}, nil
}

func NewCollectAggregate(selectors []*CollectAggregateSelector) (*CollectAggregate, error) {
	if selectors == nil {
		return nil, core.Error(core.ErrMissedArgument, "aggregate selectors")
	}

	return &CollectAggregate{selectors}, nil
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

	if clause.params.group != nil {
		grouping := clause.params.group

		for _, selector := range grouping.selectors {
			vars = append(vars, selector.variable)
		}

		if grouping.projection != nil {
			vars = append(vars, clause.params.group.projection.selector.variable)
		}

		if grouping.count != nil {
			vars = append(vars, clause.params.group.count.variable)
		}

		if grouping.aggregate != nil {
			for _, selector := range grouping.aggregate.selectors {
				vars = append(vars, selector.variable)
			}
		}
	} else if clause.params.count != nil {
		vars = append(vars, clause.params.count.variable)
	} else if clause.params.aggregate != nil {
		for _, selector := range clause.params.aggregate.selectors {
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
	)
}
