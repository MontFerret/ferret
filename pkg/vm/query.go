package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	queryDescriptorKind       = "kind"
	queryDescriptorExpression = "expression"
	queryDescriptorParams     = "params"
	queryDescriptorOptions    = "options"

	errQueryFormatUnexpected = "unexpected query format"
)

var (
	queryDescriptorKeyKind       = runtime.NewString(queryDescriptorKind)
	queryDescriptorKeyExpression = runtime.NewString(queryDescriptorExpression)
	queryDescriptorKeyParams     = runtime.NewString(queryDescriptorParams)
	queryDescriptorKeyOptions    = runtime.NewString(queryDescriptorOptions)
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

func applyQueryExists(ctx context.Context, src runtime.Value, descriptor runtime.Value) (runtime.Value, error) {
	query, err := coerceQueryDescriptor(ctx, descriptor)
	if err != nil {
		return runtime.False, err
	}

	if queryable, ok := src.(runtime.Queryable); ok {
		return queryable.QueryExists(ctx, query)
	}

	if list, ok := src.(runtime.List); ok {
		exists := runtime.False

		err := runtime.ForEach(ctx, list, func(ctx context.Context, value, _ runtime.Value) (runtime.Boolean, error) {
			queryable, ok := value.(runtime.Queryable)
			if !ok {
				return runtime.False, runtime.TypeErrorOf(value, runtime.TypeQueryable)
			}

			res, err := queryable.QueryExists(ctx, query)
			if err != nil {
				return runtime.False, err
			}

			if res {
				exists = runtime.True
				return runtime.False, nil
			}

			return runtime.True, nil
		})
		if err != nil {
			return runtime.False, err
		}

		return exists, nil
	}

	return runtime.False, runtime.TypeErrorOf(src, runtime.TypeQueryable, runtime.TypeList)
}

func applyQueryCount(ctx context.Context, src runtime.Value, descriptor runtime.Value) (runtime.Value, error) {
	query, err := coerceQueryDescriptor(ctx, descriptor)
	if err != nil {
		return runtime.ZeroInt, err
	}

	if queryable, ok := src.(runtime.Queryable); ok {
		return queryable.QueryCount(ctx, query)
	}

	if list, ok := src.(runtime.List); ok {
		total := runtime.ZeroInt

		err := runtime.ForEach(ctx, list, func(ctx context.Context, value, _ runtime.Value) (runtime.Boolean, error) {
			queryable, ok := value.(runtime.Queryable)
			if !ok {
				return runtime.False, runtime.TypeErrorOf(value, runtime.TypeQueryable)
			}

			count, err := queryable.QueryCount(ctx, query)
			if err != nil {
				return runtime.False, err
			}

			total += count

			return runtime.True, nil
		})
		if err != nil {
			return runtime.ZeroInt, err
		}

		return total, nil
	}

	return runtime.ZeroInt, runtime.TypeErrorOf(src, runtime.TypeQueryable, runtime.TypeList)
}

func applyQueryOne(ctx context.Context, src runtime.Value, descriptor runtime.Value) (runtime.Value, error) {
	query, err := coerceQueryDescriptor(ctx, descriptor)
	if err != nil {
		return runtime.None, err
	}

	if queryable, ok := src.(runtime.Queryable); ok {
		return normalizeQueryOneResult(queryable.QueryOne(ctx, query))
	}

	if list, ok := src.(runtime.List); ok {
		var result runtime.Value = runtime.None

		err := runtime.ForEach(ctx, list, func(ctx context.Context, value, _ runtime.Value) (runtime.Boolean, error) {
			queryable, ok := value.(runtime.Queryable)
			if !ok {
				return runtime.False, runtime.TypeErrorOf(value, runtime.TypeQueryable)
			}

			out, err := normalizeQueryOneResult(queryable.QueryOne(ctx, query))
			if err != nil {
				return runtime.False, err
			}
			if runtime.TypeNone.Is(out) {
				return runtime.True, nil
			}

			result = out
			return runtime.False, nil
		})
		if err != nil {
			return runtime.None, err
		}

		return result, nil
	}

	return runtime.None, runtime.TypeErrorOf(src, runtime.TypeQueryable, runtime.TypeList)
}

func normalizeQueryOneResult(value runtime.Value, err error) (runtime.Value, error) {
	if err != nil {
		return runtime.None, err
	}
	if value == nil {
		return runtime.None, nil
	}

	return value, nil
}
