package vm

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
)

type errorHandler struct {
	state *execState
}

// protected applies protected-frame unwinding policy.
func (h errorHandler) protected(err error) error {
	if err == nil {
		return nil
	}

	if h.state.unwindToProtected() {
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

	if catch, ok := h.state.vm.tryCatch(h.state.pc); ok {
		if onCatch != nil {
			onCatch()
		}

		if catch[2] >= 0 {
			h.state.pc = catch[2]
		}

		return nil
	}

	return h.protected(err)
}

func (exec *execState) wrapRuntimeError(err error) error {
	return diagnostic.WrapRuntimeError(exec.vm.program, exec.pc, err)
}

func (exec *execState) runtimeErrorFromPanic(r any) error {
	return diagnostic.RuntimeErrorFromPanic(exec.vm.program, exec.pc, r)
}

func (exec *execState) checkDivisionByZero(ctx context.Context, left, right runtime.Value) error {
	return diagnostic.CheckDivisionByZero(ctx, exec.vm.program, exec.pc, left, right)
}

func (exec *execState) checkModuloByZero(ctx context.Context, right runtime.Value) error {
	return diagnostic.CheckModuloByZero(ctx, exec.vm.program, exec.pc, right)
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
	reg := h.state.registers.Values

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
		h.state.registers.Values[dst] = val

		return nil
	}

	if optional || errors.Is(err, runtime.ErrNotFound) {
		h.state.registers.Values[dst] = runtime.None

		return nil
	}

	return err
}

func (h errorHandler) setCallResult(op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) error {
	reg := h.state.registers.Values

	if err == nil {
		reg[dst] = out

		return nil
	}

	if bytecode.IsProtectedCallOpcode(op) {
		reg[dst] = runtime.None

		return nil
	}

	if catch, ok := h.state.vm.tryCatch(h.state.pc); ok {
		reg[dst] = runtime.None

		if catch[2] >= 0 {
			h.state.pc = catch[2]
		}

		return nil
	}

	if h.state.unwindToProtected() {
		return nil
	}

	return err
}
