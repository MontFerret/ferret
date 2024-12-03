package internal

import (
	"encoding/binary"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"hash/fnv"
)

type Tuple struct {
	First  core.Value
	Second core.Value
}

func (p *Tuple) MarshalJSON() ([]byte, error) {
	panic("not supported")
}

func (p *Tuple) String() string {
	return "[Tuple]"
}

func (p *Tuple) Unwrap() interface{} {
	return [2]core.Value{p.First, p.Second}
}

func (p *Tuple) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("internal.Tuple"))
	h.Write([]byte(":"))
	h.Write([]byte("["))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, p.First.Hash())

	h.Write(bytes)
	h.Write([]byte(","))
	binary.LittleEndian.PutUint64(bytes, p.Second.Hash())
	h.Write(bytes)

	h.Write([]byte("]"))

	return h.Sum64()
}

func (p *Tuple) Copy() core.Value {
	return &Tuple{
		First:  p.First,
		Second: p.Second,
	}
}
