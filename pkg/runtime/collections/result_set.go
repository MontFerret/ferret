package collections

import (
	"encoding/binary"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"hash/fnv"
)

type ResultSet []core.Value

func (set ResultSet) Hash() uint64 {
	if set == nil {
		return 0
	}

	h := fnv.New64a()

	for i, el := range set {
		idx := make([]byte, 8)
		binary.LittleEndian.PutUint64(idx, uint64(i))
		h.Write(idx)

		val := make([]byte, 8)
		binary.LittleEndian.PutUint64(val, el.Hash())
		h.Write(val)
	}

	return h.Sum64()
}

func (set ResultSet) First() core.Value {
	if len(set) > 0 {
		return set[0]
	}

	return values.None
}

func (set ResultSet) Last() core.Value {
	if len(set) > 0 {
		return set[len(set)-1]
	}

	return values.None
}
