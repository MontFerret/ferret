package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (vm *VM) castSubscribeArgs(dst, eventName, opts runtime.Value) (runtime.Observable, runtime.String, runtime.Map, error) {
	observable, ok := dst.(runtime.Observable)

	if !ok {
		return nil, "", nil, runtime.TypeErrorOf(dst, runtime.TypeObservable)
	}

	eventNameStr, ok := eventName.(runtime.String)

	if !ok {
		return nil, "", nil, runtime.TypeErrorOf(eventName, runtime.TypeString)
	}

	var options runtime.Map

	if opts != nil && opts != runtime.None {
		m, ok := opts.(runtime.Map)

		if !ok {
			return nil, "", nil, runtime.TypeErrorOf(opts, runtime.TypeMap)
		}

		options = m
	}

	return observable, eventNameStr, options, nil
}

func (vm *VM) castDispatchArgs(
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
	var options runtime.Value = runtime.None

	if args == nil || args == runtime.None {
		return dispatcher, eventNameStr, payload, options, nil
	}

	argMap, err := runtime.CastMap(args)

	if err != nil {
		return nil, "", nil, nil, err
	}

	if val, err := argMap.Get(ctx, runtime.NewString("payload")); err == nil {
		payload = val
	}

	if val, err := argMap.Get(ctx, runtime.NewString("options")); err == nil {
		options = val
	}

	return dispatcher, eventNameStr, payload, options, nil
}
