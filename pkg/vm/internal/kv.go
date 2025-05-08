package internal

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type KeyValuePair struct {
	Key   runtime.Value
	Value runtime.Value
}

func (p *KeyValuePair) MarshalJSON() ([]byte, error) {
	return p.Value.MarshalJSON()
}

func (p *KeyValuePair) String() string {
	return "[KeyValuePair]"
}

func (p *KeyValuePair) Unwrap() interface{} {
	return [2]runtime.Value{p.Key, p.Value}
}

func (p *KeyValuePair) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("vm.KeyValuePair"))
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

func (p *KeyValuePair) Copy() runtime.Value {
	return &KeyValuePair{
		Key:   p.Key,
		Value: p.Value,
	}
}
