package arrays

import (
	"context"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func mutableListArg(ctx context.Context, value runtime.Value) (runtime.List, error) {
	if err := ctx.Err(); err != nil {
		return nil, runtime.ArgError(err, 0)
	}

	// A nil receiver can satisfy List while being unsafe to call.
	reflected := reflect.ValueOf(value)
	if !reflected.IsValid() {
		return nil, runtime.ArgTypeError(runtime.None, 0, runtime.TypeList)
	}

	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if reflected.IsNil() {
			return nil, runtime.ArgTypeError(runtime.None, 0, runtime.TypeList)
		}
	}

	return runtime.CastArg[runtime.List](value, 0)
}

func mutableListLength(ctx context.Context, list runtime.List) (runtime.Int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	size, err := list.Length(ctx)
	if err != nil {
		return 0, err
	}

	if size < 0 {
		return 0, runtime.Error(runtime.ErrInvalidOperation, "negative list length")
	}

	return size, nil
}

func mutableIndexArgs(ctx context.Context, array, index runtime.Value, allowEnd bool) (runtime.List, runtime.Int, error) {
	list, err := mutableListArg(ctx, array)
	if err != nil {
		return nil, 0, err
	}

	position, err := runtime.CastArg[runtime.Int](index, 1)
	if err != nil {
		return nil, 0, err
	}

	size, err := mutableListLength(ctx, list)
	if err != nil {
		return nil, 0, runtime.ArgError(err, 0)
	}

	if position < 0 || position > size || (!allowEnd && position == size) {
		return nil, 0, runtime.ArgError(runtime.Error(runtime.ErrInvalidOperation, "out of bounds"), 1)
	}

	return list, position, nil
}

func mutationResult(target runtime.Value, err error) (runtime.Value, error) {
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return target, nil
}
