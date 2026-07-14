package sdk

import (
	"context"
	"fmt"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SliceIterator iterates over a fixed-length view of a Go slice using an Encoder.
type SliceIterator[T any] struct {
	encoder Encoder[T]
	data    []T
	length  int
	pos     int
}

// NewSliceIterator creates an iterator using DefaultCodec.
func NewSliceIterator[T any](data []T) runtime.Iterator {
	return NewSliceIteratorWithEncoding(data, DefaultCodec[T]())
}

// NewSliceIteratorWithEncoding creates an iterator using encoder.
func NewSliceIteratorWithEncoding[T any](data []T, encoder Encoder[T]) runtime.Iterator {
	return &SliceIterator[T]{
		data:    data,
		length:  len(data),
		encoder: encoder,
	}
}

// Next encodes the next item and returns its zero-based index.
func (iterator *SliceIterator[T]) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, runtime.None, err
	}
	if iterator.pos >= iterator.length {
		return runtime.None, runtime.None, io.EOF
	}

	position := iterator.pos
	value, err := iterator.encoder.Encode(ctx, iterator.data[position])

	if err != nil {
		return runtime.None, runtime.None, fmt.Errorf("slice index %d: %w", position, err)
	}

	iterator.pos++

	return normalizeRuntimeValue(value), runtime.NewInt(position), nil
}
