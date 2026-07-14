package sdk

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MapView exposes a live Go map through its supported keyed collection capabilities.
// It intentionally does not claim the full runtime.Map contract.
type MapView[TKey comparable, TValue any] struct {
	*HostValue[map[TKey]TValue]
	keyCodec   Codec[TKey]
	valueCodec Codec[TValue]
}

// NewMapView creates a live map view using DefaultCodec for keys and values.
func NewMapView[TKey comparable, TValue any](data map[TKey]TValue) *MapView[TKey, TValue] {
	return NewMapViewWithEncoding(data, DefaultCodec[TKey](), DefaultCodec[TValue]())
}

// NewMapViewWithEncoding creates a live map view using explicit key and value codecs.
func NewMapViewWithEncoding[TKey comparable, TValue any](
	data map[TKey]TValue,
	keyCodec Codec[TKey],
	valueCodec Codec[TValue],
) *MapView[TKey, TValue] {
	return &MapView[TKey, TValue]{
		HostValue:  NewHostValue(data),
		keyCodec:   keyCodec,
		valueCodec: valueCodec,
	}
}

// Get decodes key and encodes its mapped value, returning None when absent.
func (view *MapView[TKey, TValue]) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	value, found, err := view.Lookup(ctx, key)
	if err != nil || !found {
		return runtime.None, err
	}

	return value, nil
}

// Lookup decodes key and distinguishes an absent mapping from conversion failure.
func (view *MapView[TKey, TValue]) Lookup(ctx context.Context, key runtime.Value) (runtime.Value, bool, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, false, err
	}

	decodedKey, err := view.keyCodec.Decode(ctx, key)
	if err != nil {
		return runtime.None, false, fmt.Errorf("map key: %w", err)
	}

	value, found := view.Target()[decodedKey]
	if !found {
		return runtime.None, false, nil
	}

	encoded, err := view.valueCodec.Encode(ctx, value)
	if err != nil {
		return runtime.None, false, fmt.Errorf("map value for key %v: %w", decodedKey, err)
	}

	return normalizeRuntimeValue(encoded), true, nil
}

// Set decodes key and value into the live backing map.
func (view *MapView[TKey, TValue]) Set(ctx context.Context, key, value runtime.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	decodedKey, err := view.keyCodec.Decode(ctx, key)
	if err != nil {
		return fmt.Errorf("map key: %w", err)
	}

	decodedValue, err := view.valueCodec.Decode(ctx, value)
	if err != nil {
		return fmt.Errorf("map value for key %v: %w", decodedKey, err)
	}

	data := view.Target()
	if data == nil {
		data = make(map[TKey]TValue)
		view.setTarget(data)
	}
	data[decodedKey] = decodedValue

	return nil
}

// RemoveKey decodes and removes a key from the live backing map.
func (view *MapView[TKey, TValue]) RemoveKey(ctx context.Context, key runtime.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	decodedKey, err := view.keyCodec.Decode(ctx, key)
	if err != nil {
		return fmt.Errorf("map key: %w", err)
	}

	delete(view.Target(), decodedKey)
	return nil
}

// Remove deletes the first map entry equal to value under runtime comparison semantics.
func (view *MapView[TKey, TValue]) Remove(ctx context.Context, value runtime.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	for key, item := range view.Target() {
		encoded, err := view.valueCodec.Encode(ctx, item)
		if err != nil {
			return fmt.Errorf("map value for key %v: %w", key, err)
		}

		if runtime.CompareValues(value, normalizeRuntimeValue(encoded)) == 0 {
			delete(view.Target(), key)
			return nil
		}
	}

	return nil
}

// Length returns the current map size.
func (view *MapView[TKey, TValue]) Length(ctx context.Context) (runtime.Int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	return runtime.NewInt(len(view.Target())), nil
}

// Iterate snapshots keys and converts keys and values through their codecs.
func (view *MapView[TKey, TValue]) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return NewMapIteratorWithEncoding(view.Target(), view.keyCodec, view.valueCodec), nil
}

// Copy creates a shallow view that preserves codecs, backing map, type, and identity.
func (view *MapView[TKey, TValue]) Copy() runtime.Value {
	if view == nil {
		return runtime.None
	}

	return &MapView[TKey, TValue]{
		HostValue:  view.HostValue.copyValue(),
		keyCodec:   view.keyCodec,
		valueCodec: view.valueCodec,
	}
}
