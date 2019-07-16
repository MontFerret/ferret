package clauses

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type CollectIterator struct {
	ready      bool
	values     []*core.Scope
	pos        int
	src        core.SourceMap
	params     *Collect
	dataSource collections.Iterator
}

func NewCollectIterator(
	src core.SourceMap,
	params *Collect,
	dataSource collections.Iterator,
) (*CollectIterator, error) {
	if params.group != nil {
		if params.group.selectors != nil {
			var err error
			sorters := make([]*collections.Sorter, len(params.group.selectors))

			for i, selector := range params.group.selectors {
				sorter, err := newGroupSorter(selector)

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

		if params.group.count != nil && params.group.projection != nil {
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
	}, nil
}

func newGroupSorter(selector *CollectSelector) (*collections.Sorter, error) {
	return collections.NewSorter(func(ctx context.Context, first, second *core.Scope) (int64, error) {
		f, err := selector.expression.Exec(ctx, first)

		if err != nil {
			return -1, err
		}

		s, err := selector.expression.Exec(ctx, second)

		if err != nil {
			return -1, err
		}

		return f.Compare(s), nil
	}, collections.SortDirectionAsc)
}

func (iterator *CollectIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	if !iterator.ready {
		iterator.ready = true
		groups, err := iterator.init(ctx, scope)

		if err != nil {
			return nil, err
		}

		iterator.values = groups
	}

	if len(iterator.values) > iterator.pos {
		val := iterator.values[iterator.pos]
		iterator.pos++

		return val, nil
	}

	return nil, core.ErrNoMoreData
}

func (iterator *CollectIterator) init(ctx context.Context, scope *core.Scope) ([]*core.Scope, error) {
	if iterator.params.group != nil {
		return iterator.group(ctx, scope)
	}

	if iterator.params.count != nil {
		return iterator.count(ctx, scope)
	}

	if iterator.params.aggregate != nil {
		return iterator.aggregate(ctx, scope)
	}

	return nil, core.ErrInvalidOperation
}

func (iterator *CollectIterator) group(ctx context.Context, scope *core.Scope) ([]*core.Scope, error) {
	// TODO: honestly, this code is ugly. it needs to be refactored in more chained way with much less if statements
	// slice of groups
	collected := make([]*core.Scope, 0, 10)
	// hash table of unique values
	// key is a DataSet hash
	// value is its index in result slice (collected)
	hashTable := make(map[uint64]int)

	groupSelectors := iterator.params.group.selectors
	proj := iterator.params.group.projection
	count := iterator.params.group.count
	aggr := iterator.params.group.aggregate

	// iterating over underlying data source
	for {
		// keep all defined variables in forked scopes
		// all those variables should not be available for further executions
		dataSourceScope, err := iterator.dataSource.Next(ctx, scope.Fork())

		if err != nil {
			if core.IsNoMoreData(err) {
				break
			}

			return nil, err
		}

		// this data dataSourceScope represents a data of a given iteration with values retrieved by selectors
		collectScope := scope.Fork()

		// map for calculating a hash value
		vals := make(map[string]core.Value)

		// iterate over each selector for a current data
		for _, selector := range groupSelectors {
			// execute a selector and get a value
			// e.g. COLLECT age = u.age
			value, err := selector.expression.Exec(ctx, dataSourceScope)

			if err != nil {
				return nil, err
			}

			if err := collectScope.SetVariable(selector.variable, value); err != nil {
				return nil, err
			}

			vals[selector.variable] = value
		}

		// it important to get hash value before projection and counting
		// otherwise hash value will be inaccurate
		h := values.MapHash(vals)

		_, exists := hashTable[h]

		if !exists {
			collected = append(collected, collectScope)
			hashTable[h] = len(collected) - 1

			if proj != nil {
				// create a new variable for keeping projection
				if err := collectScope.SetVariable(proj.selector.variable, values.NewArray(10)); err != nil {
					return nil, err
				}
			} else if count != nil {
				// create a new variable for keeping counter
				if err := collectScope.SetVariable(count.variable, values.ZeroInt); err != nil {
					return nil, err
				}
			} else if aggr != nil {
				// create a new variable for keeping aggregated values
				for _, selector := range aggr.selectors {
					arr := values.NewArray(len(selector.aggregators))

					for range selector.aggregators {
						arr.Push(values.None)
					}

					if err := collectScope.SetVariable(selector.variable, arr); err != nil {
						return nil, err
					}
				}
			}
		}

		if proj != nil {
			idx := hashTable[h]
			collectedScope := collected[idx]
			groupValue, err := collectedScope.GetVariable(proj.selector.variable)

			if err != nil {
				return nil, err
			}

			arr, ok := groupValue.(*values.Array)

			if !ok {
				return nil, core.TypeError(groupValue.Type(), types.Int)
			}

			value, err := proj.selector.expression.Exec(ctx, dataSourceScope)

			if err != nil {
				return nil, err
			}

			arr.Push(value)
		} else if count != nil {
			idx := hashTable[h]
			ds := collected[idx]
			groupValue, err := ds.GetVariable(count.variable)

			if err != nil {
				return nil, err
			}

			counter, ok := groupValue.(values.Int)

			if !ok {
				return nil, core.TypeError(groupValue.Type(), types.Int)
			}

			groupValue = counter + 1
			// dataSourceScope a new value
			if err := ds.UpdateVariable(count.variable, groupValue); err != nil {
				return nil, err
			}
		} else if aggr != nil {
			idx := hashTable[h]
			ds := collected[idx]

			// iterate over each selector for a current data dataSourceScope
			for _, selector := range aggr.selectors {
				sv, err := ds.GetVariable(selector.variable)

				if err != nil {
					return nil, err
				}

				vv := sv.(*values.Array)

				// execute a selector and get a value
				// e.g. AGGREGATE age = CONCAT(u.age, u.dob)
				// u.age and u.dob get executed
				for idx, exp := range selector.aggregators {
					arg, err := exp.Exec(ctx, dataSourceScope)

					if err != nil {
						return nil, err
					}

					var args *values.Array
					idx := values.NewInt(idx)

					if vv.Get(idx) == values.None {
						args = values.NewArray(10)
						vv.Set(idx, args)
					} else {
						args = vv.Get(idx).(*values.Array)
					}

					args.Push(arg)
				}
			}
		}
	}

	if aggr != nil {
		for _, iterScope := range collected {
			for _, selector := range aggr.selectors {
				sv, err := iterScope.GetVariable(selector.variable)

				if err != nil {
					return nil, err
				}

				arr := sv.(*values.Array)

				matrix := make([]core.Value, arr.Length())

				arr.ForEach(func(value core.Value, idx int) bool {
					matrix[idx] = value

					return true
				})

				reduced, err := selector.reducer(ctx, matrix...)

				if err != nil {
					return nil, err
				}

				// replace value with calculated one
				if err := iterScope.UpdateVariable(selector.variable, reduced); err != nil {
					return nil, err
				}
			}
		}
	}

	return collected, nil
}

func (iterator *CollectIterator) count(ctx context.Context, scope *core.Scope) ([]*core.Scope, error) {
	var counter int

	// iterating over underlying data source
	for {
		// keep all defined variables in forked scopes
		// all those variables should not be available for further executions
		_, err := iterator.dataSource.Next(ctx, scope.Fork())

		if err != nil {
			if core.IsNoMoreData(err) {
				break
			}

			return nil, err
		}

		counter++
	}

	cs := scope.Fork()

	if err := cs.SetVariable(iterator.params.count.variable, values.NewInt(counter)); err != nil {
		return nil, err
	}

	return []*core.Scope{cs}, nil
}

func (iterator *CollectIterator) aggregate(ctx context.Context, scope *core.Scope) ([]*core.Scope, error) {
	cs := scope.Fork()

	// matrix of aggregated expressions
	// string key of the map is a selector variable
	// value of the map is a matrix of arguments
	// e.g. x = CONCAT(arg1, arg2, argN...)
	// x is a string key where a nested array is an array of all values of argN expressions
	aggregated := make(map[string][]core.Value)
	selectors := iterator.params.aggregate.selectors

	// iterating over underlying data source
	for {
		// keep all defined variables in forked scopes
		// all those variables should not be available for further executions
		os, err := iterator.dataSource.Next(ctx, scope.Fork())

		if err != nil {
			if core.IsNoMoreData(err) {
				break
			}

			return nil, err
		}

		// iterate over each selector for a current data set
		for _, selector := range selectors {
			vv, exists := aggregated[selector.variable]

			if !exists {
				vv = make([]core.Value, len(selector.aggregators))
				aggregated[selector.variable] = vv
			}

			// execute a selector and get a value
			// e.g. AGGREGATE age = CONCAT(u.age, u.dob)
			// u.age and u.dob get executed
			for idx, exp := range selector.aggregators {
				arg, err := exp.Exec(ctx, os)

				if err != nil {
					return nil, err
				}

				var args *values.Array

				if vv[idx] == nil {
					args = values.NewArray(10)
					vv[idx] = args
				} else {
					args = vv[idx].(*values.Array)
				}

				args.Push(arg)
			}
		}
	}

	for _, selector := range selectors {
		matrix := aggregated[selector.variable]

		reduced, err := selector.reducer(ctx, matrix...)

		if err != nil {
			return nil, err
		}

		if err := cs.SetVariable(selector.variable, reduced); err != nil {
			return nil, err
		}
	}

	return []*core.Scope{cs}, nil
}
