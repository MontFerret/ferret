package runtime

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/wI2L/jettison"
)

type AggregateKind int

const (
	AggregateCount AggregateKind = iota
	AggregateSum
	AggregateMin
	AggregateMax
	AggregateAverage
)

type AggregatePlan struct {
	keys  []String
	kinds []AggregateKind
	index map[string]int
}

func NewAggregatePlan(keys []String, kinds []AggregateKind) *AggregatePlan {
	idx := make(map[string]int, len(keys))

	for i, key := range keys {
		idx[key.String()] = i
	}

	return &AggregatePlan{
		keys:  keys,
		kinds: kinds,
		index: idx,
	}
}

func (p *AggregatePlan) Keys() []String {
	return p.keys
}

func (p *AggregatePlan) Size() int {
	return len(p.keys)
}

func (p *AggregatePlan) Index(key string) (int, bool) {
	idx, ok := p.index[key]
	return idx, ok
}

func (p *AggregatePlan) KindAt(idx int) AggregateKind {
	return p.kinds[idx]
}

func (p *AggregatePlan) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(struct {
		Keys  []String        `json:"keys"`
		Kinds []AggregateKind `json:"kinds"`
	}{
		Keys:  p.keys,
		Kinds: p.kinds,
	}, jettison.NoHTMLEscaping())
}

func (p *AggregatePlan) String() string {
	data, err := p.MarshalJSON()
	if err != nil {
		return "[AggregatePlan]"
	}

	return string(data)
}

func (p *AggregatePlan) Unwrap() interface{} {
	return p
}

func (p *AggregatePlan) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("aggregate_plan:"))

	for i, key := range p.keys {
		h.Write([]byte(key))
		h.Write([]byte{0})

		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(p.kinds[i]))
		h.Write(buf[:])
	}

	return h.Sum64()
}

func (p *AggregatePlan) Copy() Value {
	keys := make([]String, len(p.keys))
	copy(keys, p.keys)

	kinds := make([]AggregateKind, len(p.kinds))
	copy(kinds, p.kinds)

	return NewAggregatePlan(keys, kinds)
}
