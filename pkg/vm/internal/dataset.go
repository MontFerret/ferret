package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type DataSet struct {
	hashmap map[uint64]bool
	// TODO: Use backend storage to support large datasets
	values runtime.List
	keyed  bool
}

func NewDataSet(distinct bool) *DataSet {
	var hashmap map[uint64]bool

	if distinct {
		hashmap = make(map[uint64]bool)
	}

	return &DataSet{
		hashmap: hashmap,
		values:  runtime.NewArray(16),
	}
}

func (ds *DataSet) Sort(ctx context.Context, direction runtime.Int) error {
	return ds.values.SortWith(ctx, func(first, second runtime.Value) int64 {
		// TODO: Brave assumption that all items are KeyValuePair
		firstKV := first.(*KeyValuePair)
		secondKV := second.(*KeyValuePair)

		comp := runtime.CompareValues(firstKV.Key, secondKV.Key)

		if direction == SortAsc {
			return comp
		}

		return -comp
	})
}

func (ds *DataSet) SortMany(ctx context.Context, directions []runtime.Int) error {
	return ds.values.SortWith(ctx, func(first, second runtime.Value) int64 {
		// TODO: Brave assumption that all items are KeyValuePair and their keys are lists
		firstKV := first.(*KeyValuePair)
		secondKV := second.(*KeyValuePair)

		firstKVKey := firstKV.Key.(runtime.List)
		secondKVKey := secondKV.Key.(runtime.List)

		for idx, direction := range directions {
			firstKey, _ := firstKVKey.Get(ctx, runtime.NewInt(idx))
			secondKey, _ := secondKVKey.Get(ctx, runtime.NewInt(idx))
			comp := runtime.CompareValues(firstKey, secondKey)

			if comp != 0 {
				if direction == SortAsc {
					return comp
				}

				return -comp
			}
		}

		return 0
	})
}

func (ds *DataSet) Get(ctx context.Context, idx runtime.Int) runtime.Value {
	val, _ := ds.values.Get(ctx, idx)

	return val
}

func (ds *DataSet) Push(ctx context.Context, item runtime.Value) {
	if ds.hashmap != nil {
		hash := item.Hash()

		_, exists := ds.hashmap[hash]

		if exists {
			return
		}

		ds.hashmap[hash] = true
	}

	_ = ds.values.Add(ctx, item)
}

func (ds *DataSet) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return ds.values.Iterate(ctx)
}

func (ds *DataSet) Length(ctx context.Context) (runtime.Int, error) {
	return ds.values.Length(ctx)
}

func (ds *DataSet) ToList() runtime.List {
	return ds.values
}

func (ds *DataSet) String() string {
	return "[DataSet]"
}

func (ds *DataSet) Unwrap() interface{} {
	return ds.values
}

func (ds *DataSet) Hash() uint64 {
	return ds.values.Hash()
}

func (ds *DataSet) Copy() runtime.Value {
	return ds.values.Copy()
}

func (ds *DataSet) MarshalJSON() ([]byte, error) {
	return ds.values.MarshalJSON()
}
