package internal

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type KV struct {
	Key   runtime.Value
	Value runtime.Value
}

func NewKV(key, value runtime.Value) *KV {
	return &KV{
		Key:   key,
		Value: value,
	}
}

func (p *KV) MarshalJSON() ([]byte, error) {
	return p.Value.MarshalJSON()
}

func (p *KV) String() string {
	return "[KV]"
}

func (p *KV) Unwrap() interface{} {
	return [2]runtime.Value{p.Key, p.Value}
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
