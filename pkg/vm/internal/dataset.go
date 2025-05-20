package internal

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type DataSet struct {
	uniqueness runtime.MapStorage
	groups     runtime.MapStorage
	values     runtime.ListStorage
	keyed      bool
}

func NewDataSet(distinct bool) *DataSet {
	var hashmap runtime.MapStorage

	if distinct {
		hashmap = runtime.NewObject()
	}

	return &DataSet{
		uniqueness: hashmap,
		values:     runtime.NewArray(16),
	}
}

func (ds *DataSet) Sort(ctx context.Context, direction runtime.Int) error {
	return ds.values.SortWith(ctx, func(first, second runtime.Value) int64 {
		firstKV, firstOk := first.(*KV)
		secondKV, secondOk := second.(*KV)

		var comp int64

		if firstOk && secondOk {
			comp = runtime.CompareValues(firstKV.Key, secondKV.Key)
		} else {
			comp = runtime.CompareValues(first, second)
		}

		if direction == SortAsc {
			return comp
		}

		return -comp
	})
}

func (ds *DataSet) SortMany(ctx context.Context, directions []runtime.Int) error {
	return ds.values.SortWith(ctx, func(first, second runtime.Value) int64 {
		firstKV, firstOk := first.(*KV)
		secondKV, secondOk := second.(*KV)

		if firstOk && secondOk {
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
		} else {
			comp := runtime.CompareValues(first, second)

			if comp != 0 {
				if directions[0] == SortAsc {
					return comp
				}

				return -comp
			}
		}

		return 0
	})
}

func (ds *DataSet) Get(ctx context.Context, idx runtime.Int) (runtime.Value, error) {
	return ds.values.Get(ctx, idx)
}

func (ds *DataSet) Add(ctx context.Context, item runtime.Value) error {
	can, err := ds.canAdd(ctx, item)

	if err != nil {
		return err
	}

	if can {
		_ = ds.values.Add(ctx, item)
	}

	return nil
}

func (ds *DataSet) AddKV(ctx context.Context, key, value runtime.Value) error {
	can, err := ds.canAdd(ctx, value)

	if err != nil {
		return err
	}

	if can {
		_ = ds.values.Add(ctx, NewKV(key, value))
		ds.keyed = true
	}

	return nil
}

func (ds *DataSet) Collect(ctx context.Context, key, value runtime.Value) error {
	return nil
}

func (ds *DataSet) Iterate(ctx context.Context) (runtime.Iterator, error) {
	iter, err := ds.values.Iterate(ctx)

	if err != nil {
		return nil, err
	}

	if !ds.keyed {
		return iter, nil
	}

	return NewKVIterator(iter), nil
}

func (ds *DataSet) Length(ctx context.Context) (runtime.Int, error) {
	return ds.values.Length(ctx)
}

func (ds *DataSet) String() string {
	return "[DataSet]"
}

func (ds *DataSet) Unwrap() interface{} {
	return ds.values
}

func (ds *DataSet) Hash() uint64 {
	return 0
}

func (ds *DataSet) Copy() runtime.Value {
	return ds
}

func (ds *DataSet) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (ds *DataSet) ToList() runtime.List {
	return ds.values
}

func (ds *DataSet) canAdd(ctx context.Context, value runtime.Value) (bool, error) {
	if ds.uniqueness == nil {
		return true, nil
	}

	hash := value.Hash()
	rnHash := runtime.Int(int64(hash))

	_, err := ds.uniqueness.Get(ctx, rnHash)

	if err != nil {
		if errors.Is(err, runtime.ErrNotFound) {
			return true, ds.uniqueness.Set(ctx, rnHash, value)
		}

		return false, err
	}

	return true, nil
}
