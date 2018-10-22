package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"sort"
)

type CollectGroupingIterator struct {
	ready      bool
	values     []collections.DataSet
	pos        int
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
		0,
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
			iterator.values = nil

			return false
		}

		iterator.values = groups
	}

	return len(iterator.values) > iterator.pos
}

func (iterator *CollectGroupingIterator) Next() (collections.DataSet, error) {
	if len(iterator.values) > iterator.pos {
		val := iterator.values[iterator.pos]
		iterator.pos++

		return val, nil
	}

	return nil, collections.ErrExhausted
}

func (iterator *CollectGroupingIterator) group() ([]collections.DataSet, error) {
	hashTable := make(map[uint64]bool)
	collected := make([]collections.DataSet, 0, 10)

	// iterating over underlying data source
	for iterator.dataSource.HasNext() {
		set, err := iterator.dataSource.Next()

		if err != nil {
			return nil, err
		}

		if len(set) == 0 {
			continue
		}

		childScope := iterator.scope.Fork()

		// populate child scope for further group operation
		// with results from an underlying source and its exposed variables
		if err := set.Apply(childScope, iterator.variables); err != nil {
			return nil, err
		}

		// represents a data of a given iteration with values retrieved by selectors
		ds := collections.NewDataSet()

		// iterate over each selector for a current data set
		for _, selector := range iterator.selectors {
			// execute a selector and get a value
			// e.g. COLLECT age = u.age
			value, err := selector.exp.Exec(iterator.ctx, childScope)

			if err != nil {
				return nil, err
			}

			ds.Set(selector.variable, value)
		}

		h := ds.Hash()

		_, exists := hashTable[h]

		if !exists {
			hashTable[h] = true
			collected = append(collected, ds)
		}
	}

	sort.SliceStable(collected, func(i, j int) bool {
		iDS := collected[i]
		jDS := collected[j]

		return iDS.Compare(jDS) < 0
	})

	return collected, nil
}
