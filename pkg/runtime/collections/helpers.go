package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func ToSlice(ctx context.Context, scope *core.Scope, iterator Iterator) ([]*core.Scope, error) {
	res := make([]*core.Scope, 0, 10)

	for {
		nextScope, err := iterator.Next(ctx, scope.Fork())

		if err != nil {
			return nil, err
		}

		if nextScope == nil {
			return res, nil
		}

		res = append(res, nextScope)
	}
}

func ToSliceCollection(ctx context.Context, iterator CollectionIterator) ([]core.Value, error) {
	res := make([]core.Value, 0, 10)

	for {
		val, _, err := iterator.Next(ctx)

		if err != nil {
			return nil, err
		}

		res = append(res, val)
	}

	return res, nil
}

func ToSliceValue(ctx context.Context, input core.Value) ([]core.Value, error) {
	switch input.Type() {
	case types.Binary,
		types.Int,
		types.Float,
		types.String,
		types.Date:

		return []core.Value{input}, nil
	case types.Array:
		arr, ok := input.(*values.Array)

		if !ok {
			return []core.Value{}, nil
		}

		slice := make([]core.Value, 0, arr.Length())

		arr.ForEach(func(value core.Value, _ int) bool {
			slice = append(slice, value)

			return true
		})

		return slice, nil
	case types.Object:
		obj, ok := input.(*values.Object)

		if !ok {
			return []core.Value{}, nil
		}

		slice := make([]core.Value, 0, obj.Length())

		obj.ForEach(func(value core.Value, _ string) bool {
			slice = append(slice, value)

			return true
		})

		return slice, nil
	default:
		iterable, ok := input.(IterableCollection)

		if ok {
			iterator, err := iterable.Iterate(ctx)

			if err != nil {
				return nil, err
			}

			return ToSliceCollection(ctx, iterator)
		}

		return []core.Value{}, nil
	}
}
