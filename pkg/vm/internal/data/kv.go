package data

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// KV represents a key-value pair where both the key and value are of type runtime.Value.
type KV struct {
	Key   runtime.Value
	Value runtime.Value
}

// NewKV creates and returns a new KV instance with the provided key and value.
func NewKV(key, value runtime.Value) *KV {
	return &KV{
		Key:   key,
		Value: value,
	}
}

func (p *KV) String() string {
	return "[KV]"
}

func (p *KV) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("vm.KV"))
	h.Write([]byte(":"))
	h.Write([]byte("["))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, p.Key.Hash())

	h.Write(bytes)
	h.Write([]byte(","))
	binary.LittleEndian.PutUint64(bytes, p.Value.Hash())
	h.Write(bytes)

	h.Write([]byte("]"))

	return h.Sum64()
}

func (p *KV) Copy() runtime.Value {
	return &KV{
		Key:   p.Key,
		Value: p.Value,
	}
}
