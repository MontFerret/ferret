package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type DataSet struct {
	values     runtime.List
	uniqueness map[uint64]bool
	keyed      bool
}

// TODO: Remove implementation of runtime.List interface. Add an unwrap opcode in the VM to unwrap the values.
// Otherwise, when it escapes to the userspace, it might cause issues with the uniqueness map.
func NewDataSet(distinct bool) runtime.List {
	var hashmap map[uint64]bool

	if distinct {
		hashmap = make(map[uint64]bool)
	}

	return &DataSet{
		uniqueness: hashmap,
		values:     runtime.NewArray(16),
	}
}

func (ds *DataSet) Sort(ctx context.Context, direction runtime.Int) error {
	return runtime.SortListWith(ctx, ds.values, func(first, second runtime.Value) int64 {
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
	return runtime.SortListWith(ctx, ds.values, func(first, second runtime.Value) int64 {
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

func (ds *DataSet) CollectGrouping(ctx context.Context) error {
	groups := make(map[string]bool)
	nextValues := runtime.NewArray(16)

	err := runtime.ForEach(ctx, ds.values, func(c context.Context, value, idx runtime.Value) (runtime.Boolean, error) {
		kv := value.(*KV)
		key, err := Stringify(c, kv.Key)

		if err != nil {
			return false, err
		}

		_, exists := groups[key]

		if !exists {
			groups[key] = true

			if err := nextValues.Add(c, kv); err != nil {
				return false, err
			}
		}

		return true, nil
	})

	if err != nil {
		return err
	}

	ds.keyed = true
	ds.values = nextValues

	return nil
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
	return ds.values.MarshalJSON()
}

func (ds *DataSet) Compare(other runtime.Value) int64 {
	return ds.values.Compare(other)
}

func (ds *DataSet) Clone(ctx context.Context) (runtime.Cloneable, error) {
	return ds.values.Clone(ctx)
}

func (ds *DataSet) Clear(ctx context.Context) error {
	return ds.values.Clear(ctx)
}

func (ds *DataSet) Set(ctx context.Context, idx runtime.Int, value runtime.Value) error {
	return ds.values.Set(ctx, idx, value)
}

func (ds *DataSet) Remove(ctx context.Context, value runtime.Value) error {
	return ds.values.Remove(ctx, value)
}

func (ds *DataSet) RemoveAt(ctx context.Context, idx runtime.Int) (runtime.Value, error) {
	return ds.values.RemoveAt(ctx, idx)
}

func (ds *DataSet) Insert(ctx context.Context, idx runtime.Int, value runtime.Value) error {
	return ds.values.Insert(ctx, idx, value)
}

func (ds *DataSet) Swap(ctx context.Context, a, b runtime.Int) error {
	return ds.values.Swap(ctx, a, b)
}

func (ds *DataSet) Find(ctx context.Context, predicate runtime.IndexedPredicate) (runtime.List, error) {
	return ds.values.Find(ctx, predicate)
}

func (ds *DataSet) FindOne(ctx context.Context, predicate runtime.IndexedPredicate) (runtime.Value, runtime.Boolean, error) {
	return ds.values.FindOne(ctx, predicate)
}

func (ds *DataSet) IndexOf(ctx context.Context, value runtime.Value) (runtime.Int, error) {
	return ds.values.IndexOf(ctx, value)
}

func (ds *DataSet) First(ctx context.Context) (runtime.Value, error) {
	return ds.values.First(ctx)
}

func (ds *DataSet) Last(ctx context.Context) (runtime.Value, error) {
	return ds.values.Last(ctx)
}

func (ds *DataSet) Slice(ctx context.Context, start, end runtime.Int) (runtime.List, error) {
	return ds.values.Slice(ctx, start, end)
}

func (ds *DataSet) ForEach(ctx context.Context, predicate runtime.IndexedPredicate) error {
	return ds.values.ForEach(ctx, predicate)
}

func (ds *DataSet) canAdd(_ context.Context, value runtime.Value) (bool, error) {
	if ds.uniqueness == nil {
		return true, nil
	}

	hash := value.Hash()

	_, exists := ds.uniqueness[hash]

	if exists {
		return false, nil
	}

	ds.uniqueness[hash] = true

	return true, nil
}
