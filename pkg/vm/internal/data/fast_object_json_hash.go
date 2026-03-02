package data

import (
	"encoding/binary"
	"hash/fnv"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (t *FastObject) MarshalJSON() ([]byte, error) {
	return json.Default.Encode(t.toMap())
}

func (t *FastObject) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(marshaled)
}

func (t *FastObject) Compare(other runtime.Value) int {
	otherObject, ok := other.(*FastObject)

	if !ok {
		if otherLike, ok := other.(runtime.ObjectLike); ok {
			return runtime.CompareTypes(t, otherLike)
		}

		return runtime.CompareTypes(t, other)
	}

	size := t.len()
	otherSize := otherObject.len()

	if size == 0 && otherSize == 0 {
		return 0
	}

	if size < otherSize {
		return -1
	}

	if size > otherSize {
		return 1
	}

	tKeys := t.keys()
	sort.Strings(tKeys)

	otherKeys := otherObject.keys()
	sort.Strings(otherKeys)

	var res int

	for i := 0; i < len(tKeys) && res == 0; i++ {
		tKey, otherKey := tKeys[i], otherKeys[i]

		if tKey == otherKey {
			tVal := t.getByKey(tKey)
			otherVal := otherObject.getByKey(otherKey)
			res = runtime.CompareValues(tVal, otherVal)
			continue
		}

		if tKey < otherKey {
			res = 1
		} else {
			res = -1
		}

		break
	}

	return res
}

func (t *FastObject) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("object:"))
	h.Write([]byte("{"))

	keys := t.keys()
	sort.Strings(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := t.getByKey(key)

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
