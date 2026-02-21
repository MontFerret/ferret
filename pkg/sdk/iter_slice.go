package sdk

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type SliceIterator[T any] struct {
	data   []T
	length int
	pos    int
}

func NewSliceIterator[T any](data []T) runtime.Iterator {
	return &SliceIterator[T]{data: data, length: len(data), pos: 0}
}

func (iter *SliceIterator[T]) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if iter.pos >= iter.length {
		return runtime.None, runtime.None, io.EOF
	}

	value := iter.data[iter.pos]
	key := runtime.NewInt(iter.pos)
	iter.pos++

	var runtimeValue runtime.Value

	if v, ok := any(value).(runtime.Value); ok {
		runtimeValue = v
	} else {
		runtimeValue = NewProxy[T](value)
	}

	return runtimeValue, key, nil
}
