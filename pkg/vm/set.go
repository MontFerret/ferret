package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
)

func (vm *VM) setIndex(ctx context.Context, target, arg, value runtime.Value) error {
	idx, err := writeIndex(arg)
	if err != nil {
		return err
	}

	writable, ok := target.(runtime.IndexWritable)
	if !ok {
		return diagnostics.MemberMutationErrorOf(target, diagnostics.MemberAccessIndex, arg)
	}

	return writable.SetAt(ctx, idx, value)
}

func (vm *VM) setKey(ctx context.Context, target, key, value runtime.Value) error {
	writable, ok := target.(runtime.KeyWritable)
	if !ok {
		return diagnostics.MemberMutationErrorOf(target, diagnostics.MemberAccessProperty, key)
	}

	return writable.Set(ctx, key, value)
}

func (vm *VM) setProperty(ctx context.Context, target, prop, value runtime.Value) error {
	switch v := prop.(type) {
	case runtime.Int:
		return vm.setIndex(ctx, target, v, value)
	case runtime.Float:
		return runtime.TypeErrorOf(v, runtime.TypeInt)
	case runtime.String:
		return vm.setKey(ctx, target, v, value)
	default:
		return vm.setKey(ctx, target, runtime.ToString(prop), value)
	}
}

func writeIndex(arg runtime.Value) (runtime.Int, error) {
	idx, ok := arg.(runtime.Int)
	if !ok {
		return 0, runtime.TypeErrorOf(arg, runtime.TypeInt)
	}

	if idx < 0 {
		return 0, runtime.Error(runtime.ErrRange, "index out of bounds")
	}

	return idx, nil
}
