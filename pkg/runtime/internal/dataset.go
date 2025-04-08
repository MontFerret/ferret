package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type DataSet struct {
	hashmap map[uint64]bool
	values  core.List
}

func NewDataSet(distinct bool) *DataSet {
	var hashmap map[uint64]bool

	if distinct {
		hashmap = make(map[uint64]bool)
	}

	return &DataSet{
		hashmap: hashmap,
		values:  core.NewArray(16),
	}
}

func (ds *DataSet) Swap(ctx context.Context, i, j core.Int) {
	ds.values.Swap(ctx, i, j)
}

func (ds *DataSet) Get(ctx context.Context, idx core.Int) core.Value {
	val, _ := ds.values.Get(ctx, idx)

	return val
}

func (ds *DataSet) Push(ctx context.Context, item core.Value) {
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

func (ds *DataSet) Iterate(ctx context.Context) (core.Iterator, error) {
	return ds.values.Iterate(ctx)
}

func (ds *DataSet) Length(ctx context.Context) (core.Int, error) {
	return ds.values.Length(ctx)
}

func (ds *DataSet) ToList() core.List {
	return ds.values
}

func (ds *DataSet) String() string {
	return "DataSet"
}

func (ds *DataSet) Unwrap() interface{} {
	return ds.values
}

func (ds *DataSet) Hash() uint64 {
	return ds.values.Hash()
}

func (ds *DataSet) Copy() core.Value {
	return ds.values.Copy()
}

func (ds *DataSet) MarshalJSON() ([]byte, error) {
	return ds.values.MarshalJSON()
}
