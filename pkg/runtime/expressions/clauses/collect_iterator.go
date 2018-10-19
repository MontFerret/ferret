package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/expressions/datasource"
)

type CollectGroupingIterator struct {
	ready      bool
	values     collections.Iterator
	src        core.SourceMap
	selectors  []*CollectSelector
	dataSource collections.Iterator
	variables  datasource.Variables
	ctx        context.Context
	scope      *core.Scope
}

func NewCollectGroupingIterator(
	src core.SourceMap,
	selectors []*CollectSelector,
	dataSource collections.Iterator,
	variables datasource.Variables,
	ctx context.Context,
	scope *core.Scope,
) *CollectGroupingIterator {
	return &CollectGroupingIterator{
		false,
		nil,
		src,
		selectors,
		dataSource,
		variables,
		ctx,
		scope,
	}
}

func (iterator *CollectGroupingIterator) HasNext() bool {
	if !iterator.ready {
		iterator.ready = true
		groups, err := iterator.group()

		if err != nil {
			iterator.values = collections.NoopIterator

			return false
		}

		iterator.values = groups
	}

	return iterator.values.HasNext()
}

func (iterator *CollectGroupingIterator) Next() (collections.ResultSet, error) {
	return iterator.values.Next()
}

func (iterator *CollectGroupingIterator) group() (collections.Iterator, error) {
	hashTable := make(map[uint64]bool)
	collectedValues := make([]core.Value, 0, 10)
	sorters := make([]*collections.Sorter, 0, 10)
	iterCounter := -1

	for iterator.dataSource.HasNext() {
		set, err := iterator.dataSource.Next()

		if err != nil {
			return nil, err
		}

		iterCounter++

		if len(set) == 0 {
			continue
		}

		childScope := iterator.scope.Fork()

		// populate child scope with values from an underlying source and its exposed variables
		if err := iterator.variables.Apply(childScope, set); err != nil {
			return nil, err
		}

		// iterate over each selector for a current data set
		for _, selector := range iterator.selectors {
			if iterCounter == 0 {
				sorter, err := collections.NewSorter(
					func(first collections.ResultSet, second collections.ResultSet) (int, error) {
						return first.First().Compare(second.First()), nil
					},
					collections.SortDirectionAsc,
				)

				if err != nil {
					return nil, err
				}

				sorters = append(sorters, sorter)
			}

			// execute a selector and get a value
			value, err := selector.exp.Exec(iterator.ctx, childScope)

			if err != nil {
				return nil, err
			}

			key := value.Hash()

			_, exists := hashTable[key]

			if !exists {
				hashTable[key] = true
				collectedValues = append(collectedValues, value)
			}
		}
	}

	return collections.NewSortIterator(
		collections.NewSliceIterator(collectedValues),
		sorters...,
	)
}
