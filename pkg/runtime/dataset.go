package runtime

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type DataSet struct {
	hashmap map[uint64]bool
	values  *values.Array
}

func NewDataSet(distinct bool) *DataSet {
	var hasmap map[uint64]bool

	if distinct {
		hasmap = make(map[uint64]bool)
	}

	return &DataSet{
		hashmap: hasmap,
		values:  values.NewArray(16),
	}
}

func (ds *DataSet) String() string {
	return ds.values.String()
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

func (ds *DataSet) ToArray() *values.Array {
	return ds.values
}
