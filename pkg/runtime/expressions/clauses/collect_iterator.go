package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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
	// key is selector variable
	// value is map of values
	hashTable := make(map[values.String]map[uint64][]core.Value)
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
		// var result collections.DataSet

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

			// use value hash as a hash for grouping
			hash := value.Hash()

			// check whether the selector is already added to the hash table
			selectorValues, exists := hashTable[selector.variable]

			if !exists {
				selectorValues = make(map[uint64][]core.Value)
				hashTable[selector.variable] = selectorValues
			}

			collectedValues, exists := selectorValues[hash]

			if !exists {
				collectedValues = make([]core.Value, 0, 10)
			}

			collectedValues = append(collectedValues, value)
		}
	}

	return collections.NoopIterator, nil
}