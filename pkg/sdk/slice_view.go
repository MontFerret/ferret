package sdk

import (
	"context"
	"fmt"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SliceView exposes a live Go slice through its supported indexed collection capabilities.
// It intentionally does not claim the full runtime.List contract.
type SliceView[T any] struct {
	*HostValue[[]T]
	codec Codec[T]
}

// NewSliceView creates a live slice view using DefaultCodec.
func NewSliceView[T any](data []T) *SliceView[T] {
	return NewSliceViewWithEncoding(data, DefaultCodec[T]())
}

// NewSliceViewWithEncoding creates a live slice view using codec.
func NewSliceViewWithEncoding[T any](data []T, codec Codec[T]) *SliceView[T] {
	return &SliceView[T]{
		HostValue: NewHostValue(data),
		codec:     codec,
	}
}

// At encodes the item at index and returns an error when index is out of bounds.
func (view *SliceView[T]) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	data := view.Target()
	if index < 0 || index >= runtime.Int(len(data)) {
		return runtime.None, runtime.Errorf(runtime.ErrRange, "slice index %d", index)
	}

	value, err := view.codec.Encode(ctx, data[index])
	if err != nil {
		return runtime.None, fmt.Errorf("slice index %d: %w", index, err)
	}

	return normalizeRuntimeValue(value), nil
}

// LookupAt safely encodes the item at index and distinguishes a missing index.
func (view *SliceView[T]) LookupAt(ctx context.Context, index runtime.Int) (runtime.Value, bool, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, false, err
	}
	if index < 0 || index >= runtime.Int(len(view.Target())) {
		return runtime.None, false, nil
	}

	value, err := view.At(ctx, index)
	return value, err == nil, err
}

// SetAt decodes value into the live backing slice at index.
func (view *SliceView[T]) SetAt(ctx context.Context, index runtime.Int, value runtime.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	data := view.Target()
	if index < 0 || index >= runtime.Int(len(data)) {
		return runtime.Errorf(runtime.ErrRange, "slice index %d", index)
	}

	decoded, err := view.codec.Decode(ctx, value)
	if err != nil {
		return fmt.Errorf("slice index %d: %w", index, err)
	}

	data[index] = decoded
	view.setTarget(data)
	return nil
}

// RemoveAt removes and encodes an item while retaining the live backing array.
func (view *SliceView[T]) RemoveAt(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	data := view.Target()
	if index < 0 {
		return runtime.None, runtime.Errorf(runtime.ErrRange, "slice index %d", index)
	}
	if index >= runtime.Int(len(data)) {
		return runtime.None, nil
	}

	removed, err := view.codec.Encode(ctx, data[index])
	if err != nil {
		return runtime.None, fmt.Errorf("slice index %d: %w", index, err)
	}

	copy(data[index:], data[index+1:])

	var zero T
	data[len(data)-1] = zero
	data = data[:len(data)-1]

	view.setTarget(data)

	return normalizeRuntimeValue(removed), nil
}

// Remove removes the first item equal to value under runtime comparison semantics.
func (view *SliceView[T]) Remove(ctx context.Context, value runtime.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	data := view.Target()
	for index, item := range data {
		encoded, err := view.codec.Encode(ctx, item)
		if err != nil {
			return fmt.Errorf("slice index %d: %w", index, err)
		}

		if runtime.CompareValues(value, normalizeRuntimeValue(encoded)) == 0 {
			_, err = view.RemoveAt(ctx, runtime.Int(index))
			return err
		}
	}

	return nil
}

// Length returns the current view length.
func (view *SliceView[T]) Length(ctx context.Context) (runtime.Int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	return runtime.NewInt(len(view.Target())), nil
}

// Iterate creates a codec-backed iterator over the current slice length.
func (view *SliceView[T]) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return NewSliceIteratorWithEncoding(view.Target(), view.codec), nil
}

// SortAsc stably sorts the live backing slice by encoded runtime values.
func (view *SliceView[T]) SortAsc(ctx context.Context) error {
	return view.sort(ctx, true)
}

// SortDesc stably sorts the live backing slice by encoded runtime values.
func (view *SliceView[T]) SortDesc(ctx context.Context) error {
	return view.sort(ctx, false)
}

// Copy creates a shallow view that preserves codec, backing storage, type, and identity.
func (view *SliceView[T]) Copy() runtime.Value {
	if view == nil {
		return runtime.None
	}

	return &SliceView[T]{
		HostValue: view.HostValue.copyValue(),
		codec:     view.codec,
	}
}

func (view *SliceView[T]) sort(ctx context.Context, ascending bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	data := view.Target()
	encoded := make([]runtime.Value, len(data))
	order := make([]int, len(data))

	for index, item := range data {
		value, err := view.codec.Encode(ctx, item)
		if err != nil {
			return fmt.Errorf("slice index %d: %w", index, err)
		}

		encoded[index] = normalizeRuntimeValue(value)
		order[index] = index
	}

	sort.SliceStable(order, func(i, j int) bool {
		comparison := runtime.CompareValues(encoded[order[i]], encoded[order[j]])
		if ascending {
			return comparison < 0
		}

		return comparison > 0
	})

	sorted := make([]T, len(data))
	for index, original := range order {
		sorted[index] = data[original]
	}
	copy(data, sorted)
	view.setTarget(data)

	return nil
}
