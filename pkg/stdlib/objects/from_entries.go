package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// FromEntries consumes an iterable of two-item, indexed, measurable pairs once.
// Keys must be strings. Duplicate keys use the last value encountered.
// @param entries {Iterable} Entries containing a String key and any value.
// @return {Map} Independent object containing cloned or copied entry values.
func FromEntries(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeIterable); err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	result := runtime.NewObject()
	index := 0
	err := runtime.ForEach(ctx, arg.(runtime.Iterable), func(ctx context.Context, entry, _ runtime.Value) (runtime.Boolean, error) {
		if err := ctx.Err(); err != nil {
			return false, err
		}

		key, value, err := readEntry(ctx, entry)
		if err != nil {
			return false, fmt.Errorf("entry %d: %w", index, err)
		}

		if err := setCopiedEntry(ctx, result, key, value); err != nil {
			return false, fmt.Errorf("entry %d: %w", index, err)
		}

		index++

		return true, nil
	})
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return result, nil
}

func readEntry(ctx context.Context, entry runtime.Value) (runtime.String, runtime.Value, error) {
	measured, ok := entry.(runtime.Measurable)
	if !ok {
		return "", nil, runtime.TypeErrorOf(entry, runtime.TypeMeasurable)
	}

	indexed, ok := entry.(runtime.IndexReadable)
	if !ok {
		return "", nil, runtime.TypeErrorOf(entry, runtime.TypeIndexReadable)
	}

	length, err := measured.Length(ctx)
	if err != nil {
		return "", nil, err
	}

	if length != 2 {
		return "", nil, runtime.Error(runtime.ErrInvalidArgument, fmt.Sprintf("expected exactly two values, got %d", length))
	}

	keyValue, err := indexed.At(ctx, 0)
	if err != nil {
		return "", nil, fmt.Errorf("key: %w", err)
	}

	key, err := runtime.CastString(keyValue)
	if err != nil {
		return "", nil, fmt.Errorf("key: %w", err)
	}

	value, err := indexed.At(ctx, 1)
	if err != nil {
		return "", nil, fmt.Errorf("value for key %q: %w", key, err)
	}

	return key, value, nil
}
