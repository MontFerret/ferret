package sdk

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Mapper represents a function that maps a value and its key to a new value of type T.
type Mapper[T any] func(ctx context.Context, value, key runtime.Value) (T, error)

// GetByKey attempts to get a value by key from the input.
// It returns an error if the input does not implement KeyReadable,
// if the key is not found, or if the found value cannot be cast to the expected type.
func GetByKey[T runtime.Value](ctx context.Context, input, key runtime.Value) (T, error) {
	keyReadable, ok := input.(runtime.KeyReadable)

	if !ok {
		var zero T

		return zero, runtime.TypeError(runtime.TypeOf(input), runtime.TypeKeyReadable)
	}

	found, err := keyReadable.Get(ctx, key)

	if err != nil {
		var zero T

		return zero, err
	}

	if found == nil || found == runtime.None {
		var zero T

		return zero, runtime.ErrNotFound
	}

	expected, ok := found.(T)

	if !ok {
		var zero T

		return zero, runtime.TypeError(runtime.TypeOf(found), runtime.TypeOf(zero))
	}

	return expected, nil
}

// GetByIndex attempts to get a value by index from the input.
// It returns an error if the input does not implement IndexReadable,
// if the index is not found, or if the found value cannot be cast to the expected type.
func GetByIndex[T runtime.Value](ctx context.Context, input runtime.Value, index runtime.Int) (T, error) {
	indexReadable, ok := input.(runtime.IndexReadable)

	if !ok {
		var zero T

		return zero, runtime.TypeError(runtime.TypeOf(input), runtime.TypeIndexReadable)
	}

	found, err := indexReadable.At(ctx, index)

	if err != nil {
		var zero T

		return zero, err
	}

	if found == nil || found == runtime.None {
		var zero T

		return zero, runtime.ErrNotFound
	}

	expected, ok := found.(T)

	if !ok {
		var zero T

		return zero, runtime.TypeError(runtime.TypeOf(found), runtime.TypeOf(zero))
	}

	return expected, nil
}

// ToSlice converts an Iterable input into a slice of type T using the provided mapper function.
// It returns an error if the input does not implement Iterable or if any mapping operation fails.
func ToSlice[T any](ctx context.Context, input runtime.Value, mapper Mapper[T]) ([]T, error) {
	iter, ok := input.(runtime.Iterable)

	if !ok {
		return nil, runtime.TypeErrorOf(input, runtime.TypeIterable)
	}

	capacity := 5

	meas, ok := input.(runtime.Measurable)

	if ok {
		res, err := meas.Length(ctx)

		if err != nil {
			return nil, err
		}

		capacity = int(res)
	}

	result := make([]T, 0, capacity)

	err := runtime.ForEach(ctx, iter, func(_ context.Context, val runtime.Value, key runtime.Value) (runtime.Boolean, error) {
		mapped, e := mapper(ctx, val, key)

		if e != nil {
			return false, e
		}

		result = append(result, mapped)

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UnwrapStrings is a helper function that takes a slice of String values and returns a slice of their underlying string representations.
func UnwrapStrings(values []runtime.String) []string {
	out := make([]string, len(values))

	for i, v := range values {
		out[i] = string(v)
	}

	return out
}

// CompareStrings compares two String values and returns an Int indicating their lexicographical order.
// It returns a negative Int if a < b, zero if a == b, and a positive Int if a > b.
func CompareStrings(a, b runtime.String) runtime.Int {
	return runtime.Int(strings.Compare(a.String(), b.String()))
}
