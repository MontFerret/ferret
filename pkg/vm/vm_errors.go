package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
)

func (vm *VM) wrapRuntimeError(err error) error {
	return diagnostic.WrapRuntimeError(vm.program, vm.pc, err)
}

func (vm *VM) runtimeErrorFromPanic(r any) error {
	return diagnostic.RuntimeErrorFromPanic(vm.program, vm.pc, r)
}

func (vm *VM) newRuntimeError(kind diagnostics.Kind, message, label, hint, note string) *RuntimeError {
	return diagnostic.NewRuntimeError(vm.program, vm.pc, kind, message, label, hint, note)
}

func (vm *VM) checkDivisionByZero(ctx context.Context, left, right runtime.Value) error {
	return diagnostic.CheckDivisionByZero(ctx, vm.program, vm.pc, left, right)
}

func (vm *VM) checkModuloByZero(ctx context.Context, right runtime.Value) error {
	return diagnostic.CheckModuloByZero(ctx, vm.program, vm.pc, right)
}
