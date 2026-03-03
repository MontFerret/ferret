package vm

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
)

// handleProtectedError applies protected-frame unwinding policy.
func (vm *VM) handleProtectedError(err error) error {
	if err == nil {
		return nil
	}

	if vm.unwindToProtected() {
		return nil
	}

	return err
}

// handleError applies catch-table then protected-frame error policy.
func (vm *VM) handleError(err error) error {
	return vm.handleErrorWithCatch(err, nil)
}

// handleErrorWithCatch applies catch-table then protected-frame error policy
// and allows a catch-specific fallback assignment/action.
func (vm *VM) handleErrorWithCatch(err error, onCatch func()) error {
	if err == nil {
		return nil
	}

	if catch, ok := vm.tryCatch(vm.pc); ok {
		if onCatch != nil {
			onCatch()
		}

		if catch[2] >= 0 {
			vm.pc = catch[2]
		}

		return nil
	}

	return vm.handleProtectedError(err)
}

func (vm *VM) wrapRuntimeError(err error) error {
	return diagnostic.WrapRuntimeError(vm.program, vm.pc, err)
}

func (vm *VM) runtimeErrorFromPanic(r any) error {
	return diagnostic.RuntimeErrorFromPanic(vm.program, vm.pc, r)
}

func (vm *VM) checkDivisionByZero(ctx context.Context, left, right runtime.Value) error {
	return diagnostic.CheckDivisionByZero(ctx, vm.program, vm.pc, left, right)
}

func (vm *VM) checkModuloByZero(ctx context.Context, right runtime.Value) error {
	return diagnostic.CheckModuloByZero(ctx, vm.program, vm.pc, right)
}

func (vm *VM) tryCatch(pos int) (bytecode.Catch, bool) {
	if vm.catchByPC != nil && pos >= 0 && pos < len(vm.catchByPC) {
		if idx := vm.catchByPC[pos]; idx >= 0 {
			return vm.program.CatchTable[idx], true
		}

		return bytecode.Catch{}, false
	}

	for _, pair := range vm.program.CatchTable {
		if pos >= pair[0] && pos <= pair[1] {
			return pair, true
		}
	}

	return bytecode.Catch{}, false
}

func (vm *VM) setOrTryCatch(dst bytecode.Operand, val runtime.Value, err error) error {
	reg := vm.registers.Values

	if err == nil {
		reg[dst] = val

		return nil
	}

	return vm.handleErrorWithCatch(err, func() {
		reg[dst] = runtime.None
	})
}

func (vm *VM) setOrOptional(dst bytecode.Operand, val runtime.Value, err error, optional bool) error {
	if err == nil {
		vm.registers.Values[dst] = val

		return nil
	}

	if optional || errors.Is(err, runtime.ErrNotFound) {
		vm.registers.Values[dst] = runtime.None

		return nil
	}

	return err
}
