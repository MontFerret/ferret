package sdk

import (
	"context"
	"fmt"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MapIterator iterates over a snapshot of Go map keys using explicit encoders.
type MapIterator[TKey comparable, TValue any] struct {
	data     map[TKey]TValue
	keyEnc   Encoder[TKey]
	valueEnc Encoder[TValue]
	keys     []TKey
	pos      int
}

// NewMapIterator creates an iterator using DefaultCodec for keys and values.
func NewMapIterator[TKey comparable, TValue any](data map[TKey]TValue) runtime.Iterator {
	return NewMapIteratorWithEncoding(data, DefaultCodec[TKey](), DefaultCodec[TValue]())
}

// NewMapIteratorWithEncoding creates an iterator using the provided key and value encoders.
func NewMapIteratorWithEncoding[TKey comparable, TValue any](
	data map[TKey]TValue,
	keyEnc Encoder[TKey],
	valueEnc Encoder[TValue],
) runtime.Iterator {
	iterator := &MapIterator[TKey, TValue]{
		data:     data,
		keys:     make([]TKey, 0, len(data)),
		keyEnc:   keyEnc,
		valueEnc: valueEnc,
	}

	for key := range data {
		iterator.keys = append(iterator.keys, key)
	}

	return iterator
}

// Next encodes the next map value and key.
func (iterator *MapIterator[TKey, TValue]) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, runtime.None, err
	}
	if iterator.pos >= len(iterator.keys) {
		return runtime.None, runtime.None, io.EOF
	}

	key := iterator.keys[iterator.pos]
	value := iterator.data[key]

	runtimeKey, err := iterator.keyEnc.Encode(ctx, key)
	if err != nil {
		return runtime.None, runtime.None, fmt.Errorf("map key at position %d: %w", iterator.pos, err)
	}

	runtimeValue, err := iterator.valueEnc.Encode(ctx, value)
	if err != nil {
		return runtime.None, runtime.None, fmt.Errorf("map value at position %d: %w", iterator.pos, err)
	}

	iterator.pos++
	return normalizeRuntimeValue(runtimeValue), normalizeRuntimeValue(runtimeKey), nil
}
