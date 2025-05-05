package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime"
	"io"
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

func (it *Iterator) HasNext(ctx context.Context) (bool, error) {
	return it.src.HasNext(ctx)
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

func (it *Iterator) Unwrap() interface{} {
	return it.src
}

func (it *Iterator) MarshalJSON() ([]byte, error) {
	panic("not supported")
}

func (it *Iterator) String() string {
	return "[Iterator]"
}

func (it *Iterator) Hash() uint64 {
	panic("not supported")
}

func (it *Iterator) Copy() runtime.Value {
	panic("not supported")
}
