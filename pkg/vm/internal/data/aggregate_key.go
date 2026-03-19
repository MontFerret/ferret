package data

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type AggregateKey struct {
	groupKey    runtime.Value
	selectorIdx int
}

func NewAggregateKey(groupKey runtime.Value, selectorIdx int) *AggregateKey {
	if groupKey == nil {
		groupKey = runtime.None
	}

	return &AggregateKey{
		groupKey:    groupKey,
		selectorIdx: selectorIdx,
	}
}

func (k *AggregateKey) GroupKey() runtime.Value {
	if k == nil || k.groupKey == nil {
		return runtime.None
	}

	return k.groupKey
}

func (k *AggregateKey) SelectorIndex() int {
	if k == nil {
		return 0
	}

	return k.selectorIdx
}

func (k *AggregateKey) String() string {
	return "[AggregateKey]"
}

func (k *AggregateKey) Hash() uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte("vm.AggregateKey"))

	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], k.GroupKey().Hash())
	_, _ = h.Write(buf[:])
	binary.LittleEndian.PutUint64(buf[:], uint64(k.SelectorIndex()))
	_, _ = h.Write(buf[:])

	return h.Sum64()
}

func (k *AggregateKey) Copy() runtime.Value {
	return &AggregateKey{
		groupKey:    k.GroupKey(),
		selectorIdx: k.SelectorIndex(),
	}
}

func (*AggregateKey) VMUntracked() {}

func DecodeAggregateKey(key runtime.Value) (runtime.Value, int, bool) {
	aggKey, ok := key.(*AggregateKey)
	if !ok {
		return nil, 0, false
	}

	return aggKey.GroupKey(), aggKey.SelectorIndex(), true
}
