package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func coerceBool(input runtime.Value) runtime.Boolean {
	if input == nil || input == runtime.None {
		return runtime.False
	}

	switch val := input.(type) {
	case runtime.Boolean:
		return val
	case runtime.String:
		return val != ""
	case runtime.Int:
		return val != 0
	case runtime.Float:
		return val != 0
	case runtime.DateTime:
		return val.IsZero() != true
	default:
		return runtime.True
	}
}

func coerceSubscribeArgs(dst, eventName, eventOpts runtime.Value) (runtime.Observable, runtime.String, runtime.Map, error) {
	observable, ok := dst.(runtime.Observable)

	if !ok {
		return nil, "", nil, runtime.TypeErrorOf(dst, runtime.TypeObservable)
	}

	eventNameStr, ok := eventName.(runtime.String)

	if !ok {
		return nil, "", nil, runtime.TypeErrorOf(eventName, runtime.TypeString)
	}

	var opts runtime.Map

	if eventOpts != nil && eventOpts != runtime.None {
		m, ok := eventOpts.(runtime.Map)

		if !ok {
			return nil, "", nil, runtime.TypeErrorOf(eventOpts, runtime.TypeMap)
		}

		opts = m
	}

	return observable, eventNameStr, opts, nil
}

func coerceDispatchArgs(
	ctx context.Context,
	target, eventName, args runtime.Value,
) (runtime.Dispatchable, runtime.String, runtime.Value, runtime.Value, error) {
	dispatcher, ok := target.(runtime.Dispatchable)

	if !ok {
		return nil, "", nil, nil, runtime.TypeErrorOf(target, runtime.TypeDispatchable)
	}

	eventNameStr, err := runtime.CastString(eventName)

	if err != nil {
		return nil, "", nil, nil, err
	}

	var payload runtime.Value = runtime.None
	var opts runtime.Value = runtime.None

	if args == nil || args == runtime.None {
		return dispatcher, eventNameStr, payload, opts, nil
	}

	argMap, err := runtime.CastMap(args)

	if err != nil {
		return nil, "", nil, nil, err
	}

	if val, err := argMap.Get(ctx, runtime.NewString("payload")); err == nil {
		payload = val
	}

	if val, err := argMap.Get(ctx, runtime.NewString("options")); err == nil {
		opts = val
	}

	return dispatcher, eventNameStr, payload, opts, nil
}

func coerceQueryDescriptor(ctx context.Context, descriptor runtime.Value) (runtime.Query, error) {
	switch value := descriptor.(type) {
	case runtime.ObjectLike:
		length, err := value.Length(ctx)
		if err != nil {
			return runtime.Query{}, err
		}
		if length != 4 {
			return runtime.Query{}, runtime.Error(runtime.ErrInvalidOperation, errQueryFormatUnexpected)
		}

		for _, key := range []runtime.Value{
			queryDescriptorKeyKind,
			queryDescriptorKeyExpression,
			queryDescriptorKeyParams,
			queryDescriptorKeyOptions,
		} {
			ok, err := value.ContainsKey(ctx, key)
			if err != nil {
				return runtime.Query{}, err
			}
			if !ok {
				return runtime.Query{}, runtime.Error(runtime.ErrInvalidOperation, errQueryFormatUnexpected)
			}
		}

		kind, err := value.Get(ctx, queryDescriptorKeyKind)
		if err != nil {
			return runtime.Query{}, err
		}

		expression, err := value.Get(ctx, queryDescriptorKeyExpression)
		if err != nil {
			return runtime.Query{}, err
		}

		params, err := value.Get(ctx, queryDescriptorKeyParams)
		if err != nil {
			return runtime.Query{}, err
		}

		options, err := value.Get(ctx, queryDescriptorKeyOptions)
		if err != nil {
			return runtime.Query{}, err
		}

		kindString, err := runtime.CastString(kind)
		if err != nil {
			return runtime.Query{}, runtime.TypeErrorOf(kind, runtime.TypeString)
		}

		expressionString, err := runtime.CastString(expression)
		if err != nil {
			return runtime.Query{}, runtime.TypeErrorOf(expression, runtime.TypeString)
		}

		return runtime.Query{
			Kind:       kindString,
			Expression: expressionString,
			Params:     params,
			Options:    options,
		}, nil
	case *runtime.Array:
		length, err := value.Length(ctx)
		if err != nil {
			return runtime.Query{}, err
		}

		if length != 4 {
			return runtime.Query{}, runtime.Error(runtime.ErrInvalidOperation, errQueryFormatUnexpected)
		}

		kindVal, err := value.At(ctx, runtime.NewInt(0))
		if err != nil {
			return runtime.Query{}, err
		}

		expressionVal, err := value.At(ctx, runtime.NewInt(1))
		if err != nil {
			return runtime.Query{}, err
		}

		paramsVal, err := value.At(ctx, runtime.NewInt(2))
		if err != nil {
			return runtime.Query{}, err
		}

		optionsVal, err := value.At(ctx, runtime.NewInt(3))
		if err != nil {
			return runtime.Query{}, err
		}

		kind, err := runtime.CastString(kindVal)
		if err != nil {
			return runtime.Query{}, runtime.TypeErrorOf(kindVal, runtime.TypeString)
		}

		expression := runtime.EmptyString
		if expressionVal != runtime.None {
			expression, err = runtime.CastString(expressionVal)
			if err != nil {
				return runtime.Query{}, runtime.TypeErrorOf(expressionVal, runtime.TypeString, runtime.TypeNone)
			}
		}

		return runtime.Query{
			Kind:       kind,
			Expression: expression,
			Params:     paramsVal,
			Options:    optionsVal,
		}, nil
	default:
		return runtime.Query{}, runtime.Error(runtime.ErrInvalidOperation, errQueryFormatUnexpected)
	}
}
