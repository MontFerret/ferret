package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type DataSet struct {
	hashmap map[uint64]bool
	values  *Array
}

func NewDataSet(distinct bool) *DataSet {
	var hashmap map[uint64]bool

	if distinct {
		hashmap = make(map[uint64]bool)
	}

	return &DataSet{
		hashmap: hashmap,
		values:  NewArray(16),
	}
}

func (ds *DataSet) Swap(i, j int) {
	ds.values.Swap(i, j)
}

func (ds *DataSet) Get(idx int) core.Value {
	return ds.values.Get(idx)
}

func (ds *DataSet) Push(item core.Value) {
	if ds.hashmap != nil {
		hash := item.Hash()

		_, exists := ds.hashmap[hash]

		if exists {
			return
		}

		ds.hashmap[hash] = true
	}

	ds.values.Push(item)
}

func (ds *DataSet) Iterate(ctx context.Context) (core.Iterator, error) {
	return ds.values.Iterate(ctx)
}

func (ds *DataSet) Length() int {
	return ds.values.Length()
}

func (ds *DataSet) ToArray() *Array {
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

func (ds *DataSet) Copy() core.Value {
	return ds.values.Copy()
}

func (ds *DataSet) MarshalJSON() ([]byte, error) {
	return ds.values.MarshalJSON()
}
