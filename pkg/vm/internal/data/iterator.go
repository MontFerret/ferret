package data

import (
	"context"
	"hash/fnv"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Iterator struct {
	src   runtime.Iterator
	value runtime.Value
	key   runtime.Value
}

var NoopIter = NewIterator(&noopIter{})

func NewIterator(src runtime.Iterator) *Iterator {
	return &Iterator{src, runtime.None, runtime.None}
}

func (it *Iterator) Next(ctx context.Context) error {
	val, key, err := it.src.Next(ctx)

	if err != nil {
		return err
	}

	it.value = val
	it.key = key

	return nil
}

func (it *Iterator) Value() runtime.Value {
	return it.value
}

func (it *Iterator) Key() runtime.Value {
	return it.key
}

func (it *Iterator) Close() error {
	if closable, ok := it.src.(io.Closer); ok {
		return closable.Close()
	}

	return nil
}

func (it *Iterator) MarshalJSON() ([]byte, error) {
	return nil, runtime.Errorf(runtime.ErrUnexpected, "iterator does not support JSON encoding")
}

func (it *Iterator) String() string {
	return "[Iterator]"
}

func (it *Iterator) Hash() uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte("vm.iterator"))

	return hasher.Sum64()
}

func (it *Iterator) Copy() runtime.Value {
	return it
}
