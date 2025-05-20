package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime"
)

type KVIterator struct {
	source runtime.Iterator
}

func NewKVIterator(source runtime.Iterator) runtime.Iterator {
	return &KVIterator{
		source: source,
	}
}

func (iter *KVIterator) HasNext(ctx context.Context) (bool, error) {
	return iter.source.HasNext(ctx)
}

func (iter *KVIterator) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	value, key, err := iter.source.Next(ctx)

	if err != nil {
		return runtime.None, runtime.None, err
	}

	if value == nil {
		return runtime.None, runtime.None, nil
	}

	kv, ok := value.(*KV)

	if ok {
		return kv.Value, kv.Key, nil
	}

	return value, key, nil
}
