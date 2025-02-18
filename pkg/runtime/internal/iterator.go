package internal

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Iterator struct {
	src   core.Iterator
	value core.Value
	key   core.Value
}

func NewIterator(src core.Iterator) *Iterator {
	return &Iterator{src, values.None, values.None}
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

func (it *Iterator) Value() core.Value {
	return it.value
}

func (it *Iterator) Key() core.Value {
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

func (it *Iterator) Copy() core.Value {
	panic("not supported")
}
