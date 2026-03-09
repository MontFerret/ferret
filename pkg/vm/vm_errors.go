package vm

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
)

type errorHandler struct {
	vm *VM
}

// protected applies protected-frame unwinding policy.
func (h errorHandler) protected(err error) error {
	if err == nil {
		return nil
	}

	if h.vm.unwindToProtected() {
		return nil
	}

	return err
}

// handle applies catch-table then protected-frame error policy.
func (h errorHandler) handle(err error) error {
	return h.handleWithCatch(err, nil)
}

// handleWithCatch applies catch-table then protected-frame error policy
// and allows a catch-specific fallback assignment/action.
func (h errorHandler) handleWithCatch(err error, onCatch func()) error {
	if err == nil {
		return nil
	}

	if catch, ok := h.vm.tryCatch(h.vm.pc); ok {
		if onCatch != nil {
			onCatch()
		}

		if catch[2] >= 0 {
			h.vm.pc = catch[2]
		}

		return nil
	}

	return h.protected(err)
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

func (h errorHandler) setOrCatch(dst bytecode.Operand, val runtime.Value, err error) error {
	reg := h.vm.registers.Values

	if err == nil {
		reg[dst] = val

		return nil
	}

	return h.handleWithCatch(err, func() {
		reg[dst] = runtime.None
	})
}

func (h errorHandler) setOptional(dst bytecode.Operand, val runtime.Value, err error, optional bool) error {
	if err == nil {
		h.vm.registers.Values[dst] = val

		return nil
	}

	if optional || errors.Is(err, runtime.ErrNotFound) {
		h.vm.registers.Values[dst] = runtime.None

		return nil
	}

	return err
}

func (h errorHandler) setCallResult(op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) error {
	reg := h.vm.registers.Values

	if err == nil {
		reg[dst] = out

		return nil
	}

	if bytecode.IsProtectedCallOpcode(op) {
		reg[dst] = runtime.None

		return nil
	}

	if catch, ok := h.vm.tryCatch(h.vm.pc); ok {
		reg[dst] = runtime.None

		if catch[2] >= 0 {
			h.vm.pc = catch[2]
		}

		return nil
	}

	if h.vm.unwindToProtected() {
		return nil
	}

	return err
}
