package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
)

func (vm *VM) deleteKey(ctx context.Context, target, key runtime.Value) error {
	removable, ok := target.(runtime.KeyRemovable)
	if !ok {
		return diagnostics.MemberDeletionErrorOf(target, diagnostics.MemberAccessProperty, key)
	}

	return removable.RemoveKey(ctx, key)
}

func (vm *VM) deleteProperty(ctx context.Context, target, prop runtime.Value) error {
	switch v := prop.(type) {
	case runtime.String:
		return vm.deleteKey(ctx, target, v)
	default:
		return vm.deleteKey(ctx, target, runtime.ToString(prop))
	}
}
