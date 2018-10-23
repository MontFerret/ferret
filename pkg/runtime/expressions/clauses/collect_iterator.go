package clauses

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type CollectIterator struct {
	ready      bool
	values     []collections.DataSet
	pos        int
	src        core.SourceMap
	params     *Collect
	dataSource collections.Iterator
	variables  collections.Variables
	ctx        context.Context
	scope      *core.Scope
}

func NewCollectIterator(
	src core.SourceMap,
	params *Collect,
	dataSource collections.Iterator,
	variables collections.Variables,
	ctx context.Context,
	scope *core.Scope,
) (*CollectIterator, error) {
	if params.Grouping != nil {
		if params.Grouping.Selectors != nil {
			var err error
			sorters := make([]*collections.Sorter, len(params.Grouping.Selectors))

			for i, selector := range params.Grouping.Selectors {
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

		if params.Grouping.Count != nil && params.Grouping.Projection != nil {
			return nil, core.Error(core.ErrInvalidArgumentNumber, "counter and projection cannot be used together")
		}
	}

	return &CollectIterator{
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

func (iterator *CollectIterator) HasNext() bool {
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

func (iterator *CollectIterator) Next() (collections.DataSet, error) {
	if len(iterator.values) > iterator.pos {
		val := iterator.values[iterator.pos]
		iterator.pos++

		return val, nil
	}

	return nil, collections.ErrExhausted
}

func (iterator *CollectIterator) init() ([]collections.DataSet, error) {
	if iterator.params.Grouping != nil {
		return iterator.group()
	}

	return nil, core.ErrNotImplemented
}

func (iterator *CollectIterator) group() ([]collections.DataSet, error) {
	// slice of groups
	collected := make([]collections.DataSet, 0, 10)
	// hash table of unique values
	// key is a DataSet hash
	// value is its index in result slice (collected)
	hashTable := make(map[uint64]int)
	ctx := iterator.ctx

	groupSelectors := iterator.params.Grouping.Selectors
	proj := iterator.params.Grouping.Projection
	count := iterator.params.Grouping.Count

	// iterating over underlying data source
	for iterator.dataSource.HasNext() {
		set, err := iterator.dataSource.Next()

		if err != nil {
			return nil, err
		}

		if len(set) == 0 {
			continue
		}

		// creating a new scope for all further operations
		childScope := iterator.scope.Fork()

		// populate the new scope with results from an underlying source and its exposed variables
		if err := set.Apply(childScope, iterator.variables); err != nil {
			return nil, err
		}

		// this data set represents a data of a given iteration with values retrieved by selectors
		ds := collections.NewDataSet()

		// iterate over each selector for a current data set
		for _, selector := range groupSelectors {
			// execute a selector and get a value
			// e.g. COLLECT age = u.age
			value, err := selector.expression.Exec(ctx, childScope)

			if err != nil {
				return nil, err
			}

			ds.Set(selector.variable, value)
		}

		// it important to get hash value before projection and counting
		// otherwise hash value will be inaccurate
		h := ds.Hash()

		_, exists := hashTable[h]

		if !exists {
			collected = append(collected, ds)
			hashTable[h] = len(collected) - 1

			if proj != nil {
				// create a new variable for keeping projection
				ds.Set(proj.selector.variable, values.NewArray(10))
			} else if count != nil {
				// create a new variable for keeping counter
				ds.Set(count.variable, values.ZeroInt)
			}
		}

		if proj != nil {
			idx := hashTable[h]
			ds := collected[idx]
			groupValue := ds.Get(proj.selector.variable)

			arr, ok := groupValue.(*values.Array)

			if !ok {
				return nil, core.TypeError(groupValue.Type(), core.IntType)
			}

			value, err := proj.selector.expression.Exec(ctx, childScope)

			if err != nil {
				return nil, err
			}

			arr.Push(value)
		} else if count != nil {
			idx := hashTable[h]
			ds := collected[idx]
			groupValue := ds.Get(count.variable)

			counter, ok := groupValue.(values.Int)

			if !ok {
				return nil, core.TypeError(groupValue.Type(), core.IntType)
			}

			groupValue = counter + 1
			// set a new value
			ds.Set(count.variable, groupValue)
		}
	}

	return collected, nil
}
