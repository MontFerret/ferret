package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
)

func (vm *VM) deleteKey(ctx context.Context, target, key runtime.Value) error {
	removable, ok := target.(runtime.KeyRemovable)
	if ok {
		return removable.RemoveKey(ctx, key)
	}

	switch v := key.(type) {
	case runtime.Int:
		return vm.deleteIndex(ctx, target, v)
	case runtime.Float:
		return runtime.TypeErrorOf(v, runtime.TypeInt)
	default:
		return diagnostics.MemberDeletionErrorOf(target, diagnostics.MemberAccessProperty, key)
	}
}

func (vm *VM) deleteIndex(ctx context.Context, target, arg runtime.Value) error {
	idx, err := writeIndex(arg)
	if err != nil {
		return err
	}

	removable, ok := target.(runtime.IndexRemovable)
	if !ok {
		return diagnostics.MemberDeletionErrorOf(target, diagnostics.MemberAccessIndex, arg)
	}

	_, err = removable.RemoveAt(ctx, idx)

	return err
}

func (vm *VM) deleteProperty(ctx context.Context, target, prop runtime.Value) error {
	switch v := prop.(type) {
	case runtime.Int:
		return vm.deleteIndex(ctx, target, v)
	case runtime.Float:
		return runtime.TypeErrorOf(v, runtime.TypeInt)
	case runtime.String:
		return vm.deleteKey(ctx, target, v)
	default:
		return vm.deleteKey(ctx, target, runtime.ToString(prop))
	}
}
