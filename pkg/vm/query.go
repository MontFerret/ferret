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

func ApplyQuery(ctx context.Context, src runtime.Value, descriptor runtime.Value) (runtime.Value, error) {
	query, err := parseQueryDescriptor(ctx, descriptor)
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

func parseQueryDescriptor(ctx context.Context, descriptor runtime.Value) (runtime.Query, error) {
	switch value := descriptor.(type) {
	case runtime.ObjectLike:
		kind, err := value.Get(ctx, queryDescriptorKeyKind)
		if err != nil {
			return runtime.Query{}, err
		}

		payload, err := value.Get(ctx, queryDescriptorKeyPayload)
		if err != nil {
			return runtime.Query{}, err
		}

		options, err := value.Get(ctx, queryDescriptorKeyOptions)
		if err != nil {
			return runtime.Query{}, err
		}

		return runtime.Query{
			Kind:    runtime.CastOr[runtime.String](kind, runtime.EmptyString),
			Payload: runtime.CastOr[runtime.String](payload, runtime.EmptyString),
			Options: options,
		}, nil
	case *runtime.Array:
		length, err := value.Length(ctx)
		if err != nil {
			return runtime.Query{}, err
		}

		if length != 3 {
			return runtime.Query{}, runtime.Error(runtime.ErrInvalidOperation, errQueryFormatUnexpected)
		}

		kindVal, err := value.At(ctx, runtime.NewInt(0))
		if err != nil {
			return runtime.Query{}, err
		}

		payloadVal, err := value.At(ctx, runtime.NewInt(1))
		if err != nil {
			return runtime.Query{}, err
		}

		optionsVal, err := value.At(ctx, runtime.NewInt(2))
		if err != nil {
			return runtime.Query{}, err
		}

		kind, err := runtime.CastString(kindVal)
		if err != nil {
			return runtime.Query{}, runtime.TypeErrorOf(kindVal, runtime.TypeString)
		}

		payload := runtime.EmptyString
		if payloadVal != runtime.None {
			payload, err = runtime.CastString(payloadVal)
			if err != nil {
				return runtime.Query{}, runtime.TypeErrorOf(payloadVal, runtime.TypeString, runtime.TypeNone)
			}
		}

		return runtime.Query{
			Kind:    kind,
			Payload: payload,
			Options: optionsVal,
		}, nil
	default:
		return runtime.Query{}, runtime.Error(runtime.ErrInvalidOperation, errQueryFormatUnexpected)
	}
}
