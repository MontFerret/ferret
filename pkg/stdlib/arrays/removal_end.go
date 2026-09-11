package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func removeEndMutable(ctx context.Context, array runtime.Value, last bool) (runtime.Value, error) {
	list, err := mutableListArg(ctx, array)
	if err != nil {
		return runtime.None, err
	}

	size, err := mutableListLength(ctx, list)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	if size == 0 {
		return runtime.None, nil
	}

	index := runtime.ZeroInt
	if last {
		index = size - 1
	}

	return removeMutableElement(ctx, list, index)
}

func removeMutableElement(ctx context.Context, list runtime.List, index runtime.Int) (runtime.Value, error) {
	value, err := list.RemoveAt(ctx, index)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	if value == nil {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidOperation, "list removal returned nil"), 0)
	}

	return value, nil
}
