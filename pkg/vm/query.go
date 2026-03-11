package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	queryDescriptorKind    = "kind"
	queryDescriptorPayload = "payload"
	queryDescriptorOptions = "options"

	errQueryFormatUnexpected = "unexpected query format"
)

var (
	queryDescriptorKeyKind    = runtime.NewString(queryDescriptorKind)
	queryDescriptorKeyPayload = runtime.NewString(queryDescriptorPayload)
	queryDescriptorKeyOptions = runtime.NewString(queryDescriptorOptions)
)

func applyQuery(ctx context.Context, src runtime.Value, descriptor runtime.Value) (runtime.Value, error) {
	query, err := coerceQueryDescriptor(ctx, descriptor)
	if err != nil {
		return runtime.None, err
	}

	if queryable, ok := src.(runtime.Queryable); ok {
		res, err := queryable.Query(ctx, query)
		if err == nil && res == nil {
			res = runtime.NewArray(0)
		}

		return res, err
	}

	if list, ok := src.(runtime.List); ok {
		out := runtime.NewArray(0)

		err := runtime.ForEach(ctx, list, func(ctx context.Context, value, _ runtime.Value) (runtime.Boolean, error) {
			queryable, ok := value.(runtime.Queryable)
			if !ok {
				return runtime.False, runtime.TypeErrorOf(value, runtime.TypeQueryable)
			}

			res, err := queryable.Query(ctx, query)
			if err != nil {
				return runtime.False, err
			}
			if res == nil {
				return runtime.True, nil
			}

			if err := runtime.ForEach(ctx, res, func(ctx context.Context, item, _ runtime.Value) (runtime.Boolean, error) {
				if err := out.Append(ctx, item); err != nil {
					return runtime.False, err
				}

				return runtime.True, nil
			}); err != nil {
				return runtime.False, err
			}

			return runtime.True, nil
		})
		if err != nil {
			return runtime.None, err
		}

		return out, nil
	}

	return runtime.None, runtime.TypeErrorOf(src, runtime.TypeQueryable, runtime.TypeList)
}
