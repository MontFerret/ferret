package core

import "context"

type key int

const paramsKey key = 0

func ParamsWith(ctx context.Context, params map[string]Value) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}

func ParamsFrom(ctx context.Context) (map[string]Value, error) {
	val := ctx.Value(paramsKey)

	param, ok := val.(map[string]Value)

	if !ok {
		return nil, Error(ErrNotFound, "parameters")
	}

	return param, nil
}

func ParamFrom(ctx context.Context, name string) (Value, error) {
	params, err := ParamsFrom(ctx)

	if err != nil {
		return nil, err
	}

	param, found := params[name]

	if !found {
		return nil, Error(ErrNotFound, "parameter."+name)
	}

	return param, nil
}
