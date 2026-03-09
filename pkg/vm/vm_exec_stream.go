package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func (vm *VM) execStreamOp(
	ctx context.Context,
	dst, src1, src2 bytecode.Operand,
	reg []runtime.Value,
) error {
	observable, eventName, options, err := vm.castSubscribeArgs(reg[dst], reg[src1], reg[src2])

	if err != nil {
		return vm.errors.handle(err)
	}

	stream, err := observable.Subscribe(ctx, runtime.Subscription{
		EventName: eventName,
		Options:   options,
	})

	if err != nil {
		return vm.errors.handle(err)
	}

	reg[dst] = data.NewStreamValue(stream)

	return nil
}

func (vm *VM) execStreamIterOp(
	_ context.Context,
	dst, src1, src2 bytecode.Operand,
	reg []runtime.Value,
) error {
	stream := reg[src1].(*data.StreamValue)

	var timeout runtime.Int

	if reg[src2] != nil && reg[src2] != runtime.None {
		t, err := runtime.CastInt(reg[src2])

		if err != nil {
			return vm.errors.handle(err)
		}

		timeout = t
	}

	reg[dst] = stream.Iterate(timeout)

	return nil
}

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
