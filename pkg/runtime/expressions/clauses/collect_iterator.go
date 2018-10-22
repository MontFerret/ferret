package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type CollectGroupingIterator struct {
	ready      bool
	values     []collections.DataSet
	pos        int
	src        core.SourceMap
	params     *CollectParams
	dataSource collections.Iterator
	variables  collections.Variables
	ctx        context.Context
	scope      *core.Scope
}

func NewCollectGroupingIterator(
	src core.SourceMap,
	params *CollectParams,
	dataSource collections.Iterator,
	variables collections.Variables,
	ctx context.Context,
	scope *core.Scope,
) (*CollectGroupingIterator, error) {
	if params.Grouping != nil {
		var err error
		sorters := make([]*collections.Sorter, len(params.Grouping))

		for i, selector := range params.Grouping {
			sorter, err := newGroupSorter(ctx, scope, variables, selector)

			if err != nil {
				return nil, err
			}

			sorters[i] = sorter
		}

		dataSource, err = collections.NewSortIterator(dataSource, sorters...)

		if err != nil {
			return nil, err
		}
	}

	return &CollectGroupingIterator{
		false,
		nil,
		0,
		src,
		params,
		dataSource,
		variables,
		ctx,
		scope,
	}, nil
}

func newGroupSorter(ctx context.Context, scope *core.Scope, variables collections.Variables, selector *CollectSelector) (*collections.Sorter, error) {
	return collections.NewSorter(func(first collections.DataSet, second collections.DataSet) (int, error) {
		scope1 := scope.Fork()
		first.Apply(scope1, variables)

		f, err := selector.expression.Exec(ctx, scope1)

		if err != nil {
			return -1, err
		}

		scope2 := scope.Fork()
		second.Apply(scope2, variables)

		s, err := selector.expression.Exec(ctx, scope2)

		if err != nil {
			return -1, err
		}

		return f.Compare(s), nil
	}, collections.SortDirectionAsc)
}

func (iterator *CollectGroupingIterator) HasNext() bool {
	if !iterator.ready {
		iterator.ready = true
		groups, err := iterator.init()

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

func (iterator *CollectGroupingIterator) init() ([]collections.DataSet, error) {
	hashTable := make(map[uint64]core.Value)
	collected := make([]collections.DataSet, 0, 10)
	ctx := iterator.ctx

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

		// populate child scope for further init operation
		// with results from an underlying source and its exposed variables
		if err := set.Apply(childScope, iterator.variables); err != nil {
			return nil, err
		}

		// represents a data of a given iteration with values retrieved by selectors
		ds := collections.NewDataSet()
		grouping := iterator.params.Grouping
		proj := iterator.params.Projection
		applyProjection := proj != nil

		// grouping can be omitted
		if grouping != nil {
			// iterate over each selector for a current data set
			for _, selector := range grouping {
				// execute a selector and get a value
				// e.g. COLLECT age = u.age
				value, err := selector.expression.Exec(ctx, childScope)

				if err != nil {
					return nil, err
				}

				ds.Set(selector.variable, value)
			}

			h := ds.Hash()

			groupValue, exists := hashTable[h]

			if !exists {
				if !applyProjection {
					groupValue = values.True
				} else {
					groupValue = values.NewArray(10)
				}

				hashTable[h] = groupValue
				collected = append(collected, ds)
			}

			if applyProjection {
				value, err := proj.selector.expression.Exec(ctx, childScope)

				if err != nil {
					return nil, err
				}

				(groupValue.(*values.Array)).Push(value)

				ds.Set(proj.selector.variable, groupValue)
			}
		}
	}

	return collected, nil
}
