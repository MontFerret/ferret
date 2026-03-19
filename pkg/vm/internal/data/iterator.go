package data

import (
	"context"
	"hash/fnv"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	IteratorState interface {
		runtime.Value
		Next(context.Context) error
		Value() runtime.Value
		Key() runtime.Value
	}

	Iterator struct {
		iteratorState
	}

	ClosableIterator struct {
		iteratorState
		closer io.Closer
	}

	iteratorState struct {
		src   runtime.Iterator
		value runtime.Value
		key   runtime.Value
	}
)

var NoopIter = NewIterator(&noopIter{})

func NewIterator(src runtime.Iterator) *Iterator {
	return &Iterator{
		iteratorState: newIteratorState(src),
	}
}

func NewClosableIterator(src runtime.Iterator) *ClosableIterator {
	closer, _ := src.(io.Closer)

	return &ClosableIterator{
		iteratorState: newIteratorState(src),
		closer:        closer,
	}
}

func WrapIterator(src runtime.Iterator) IteratorState {
	if _, ok := src.(io.Closer); ok {
		return NewClosableIterator(src)
	}

	return NewIterator(src)
}

func newIteratorState(src runtime.Iterator) iteratorState {
	return iteratorState{
		src:   src,
		value: runtime.None,
		key:   runtime.None,
	}
}

func (it *iteratorState) Next(ctx context.Context) error {
	val, key, err := it.src.Next(ctx)

	if err != nil {
		return err
	}

	it.value = val
	it.key = key

	return nil
}

func (it *iteratorState) Value() runtime.Value {
	return it.value
}

func (it *iteratorState) Key() runtime.Value {
	return it.key
}

func (*Iterator) VMUntracked() {}

func (it *Iterator) MarshalJSON() ([]byte, error) {
	return nil, runtime.Errorf(runtime.ErrUnexpected, "iterator does not support JSON encoding")
}

func (it *ClosableIterator) MarshalJSON() ([]byte, error) {
	return nil, runtime.Errorf(runtime.ErrUnexpected, "iterator does not support JSON encoding")
}

func (it *Iterator) String() string {
	return "[Iterator]"
}

func (it *ClosableIterator) String() string {
	return "[Iterator]"
}

func (it *Iterator) Hash() uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte("vm.iterator"))

	return hasher.Sum64()
}

func (it *ClosableIterator) Hash() uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte("vm.iterator"))

	return hasher.Sum64()
}

func (it *Iterator) Copy() runtime.Value {
	return it
}

func (it *ClosableIterator) Copy() runtime.Value {
	return it
}

func (it *ClosableIterator) Close() error {
	if it.closer == nil {
		return nil
	}

	return it.closer.Close()
}
