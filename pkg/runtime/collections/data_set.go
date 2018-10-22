package collections

import (
	"encoding/binary"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"hash/fnv"
	"sort"
)

type DataSet map[string]core.Value

func NewDataSet() DataSet {
	return make(DataSet)
}

func (ds DataSet) Apply(scope *core.Scope, variables Variables) error {
	if err := ValidateDataSet(ds, variables); err != nil {
		return err
	}

	for _, variable := range variables {
		if variable != "" {
			value, found := ds[variable]

			if !found {
				return core.Errorf(core.ErrNotFound, "variable not found in a given data set: %s", variable)
			}

			scope.SetVariable(variable, value)
		}
	}

	return nil
}

func (ds DataSet) Set(key string, value core.Value) {
	ds[key] = value
}

func (ds DataSet) Get(key string) core.Value {
	val, found := ds[key]

	if found {
		return val
	}

	return values.None
}

func (ds DataSet) Hash() uint64 {
	h := fnv.New64a()

	keys := make([]string, 0, len(ds))

	for key := range ds {
		keys = append(keys, key)
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	h.Write([]byte("{"))

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := ds[key]

		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, el.Hash())

		h.Write(bytes)

		if idx != endIndex {
			h.Write([]byte(","))
		}
	}

	h.Write([]byte("}"))

	return h.Sum64()
}

func (ds DataSet) Compare(other DataSet) int {
	if len(ds) > len(ds) {
		return 1
	}

	if len(ds) < len(ds) {
		return -1
	}

	var res = 0

	for key, otherVal := range other {
		res = -1

		if val, exists := ds[key]; exists {
			res = val.Compare(otherVal)
		}
	}

	return res
}
