package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type CollectGroupingIterator struct {
	ready      bool
	result     collections.Iterator
	src        core.SourceMap
	selectors  []*CollectSelector
	dataSource collections.Iterator
	variables  collections.Variables
	ctx        context.Context
	scope      *core.Scope
}

func NewCollectGroupingIterator(
	src core.SourceMap,
	selectors []*CollectSelector,
	dataSource collections.Iterator,
	variables collections.Variables,
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
			iterator.result = collections.NoopIterator

			return false
		}

		iterator.result = groups
	}

	return iterator.result.HasNext()
}

func (iterator *CollectGroupingIterator) Next() (collections.DataSet, error) {
	return iterator.result.Next()
}

func (iterator *CollectGroupingIterator) group() (collections.Iterator, error) {
	hashTable := make(map[uint64]bool)
	collectedValues := make([]collections.DataSet, 0, 10)
	// sorters := make([]*collections.Sorter, 0, 10)
	iterCounter := -1

	// iterating over underlying data source
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

		// populate child scope for further group operation
		// with results from an underlying source and its exposed variables
		if err := set.Apply(childScope, iterator.variables); err != nil {
			return nil, err
		}

		// result list of the current iteration
		// if there are no unique values, it will be nil
		var result collections.DataSet

		// iterate over each selector for a current data set
		for _, selector := range iterator.selectors {
			//if iterCounter == 0 {
			//	sorter, err := collections.NewSorter(
			//		func(first collections.ResultSet, second collections.ResultSet) (int, error) {
			//			return first.First().Compare(second.First()), nil
			//		},
			//		collections.SortDirectionAsc,
			//	)
			//
			//	if err != nil {
			//		return nil, err
			//	}
			//
			//	sorters = append(sorters, sorter)
			//}

			// execute a selector and get a value
			// e.g. COLLECT age = u.age
			value, err := selector.exp.Exec(iterator.ctx, childScope)

			if err != nil {
				return nil, err
			}

			// use value hash as a key for grouping
			key := value.Hash()

			// check whether the value already is added to the hash table
			_, exists := hashTable[key]

			if !exists {
				hashTable[key] = true

				if result == nil {
					result = make(collections.DataSet)
				}

				// result[selector.variable] =
			}
		}

		// put the data set of the current iteration to the final list
		if result != nil {
			collectedValues = append(collectedValues, result)
		}
	}

	return collections.NoopIterator, nil
}
