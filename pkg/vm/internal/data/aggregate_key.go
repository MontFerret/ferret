package data

import (
	"context"
	"encoding/binary"
	"hash/fnv"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
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

func (*AggregateKey) VMDefinitelyNonOwning() {}

func DecodeAggregateKey(ctx context.Context, key runtime.Value) (runtime.Value, int, bool, error) {
	if aggKey, ok := key.(*AggregateKey); ok {
		return aggKey.GroupKey(), aggKey.SelectorIndex(), true, nil
	}

	list, ok := key.(runtime.List)
	if !ok {
		return nil, 0, false, nil
	}

	length, err := list.Length(ctx)
	if err != nil {
		return nil, 0, false, err
	}

	if length != 3 {
		return nil, 0, false, nil
	}

	marker, err := list.At(ctx, 0)
	if err != nil {
		return nil, 0, false, err
	}

	if marker != bytecode.AggregateKeyMarker {
		return nil, 0, false, nil
	}

	groupKey, err := list.At(ctx, 1)
	if err != nil {
		return nil, 0, false, err
	}

	idxVal, err := list.At(ctx, 2)
	if err != nil {
		return nil, 0, false, err
	}

	idx, ok := idxVal.(runtime.Int)
	if !ok {
		return nil, 0, false, runtime.Errorf(runtime.ErrInvalidArgument, "aggregate selector index invalid")
	}

	return groupKey, int(idx), true, nil
}
