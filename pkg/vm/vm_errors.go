package vm

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
)

// handleProtectedError applies protected-frame unwinding policy.
func (s *execState) handleProtectedError(err error) error {
	if err == nil {
		return nil
	}

	if s.unwindToProtected() {
		return nil
	}

	return err
}

// handleError applies catch-table then protected-frame error policy.
func (s *execState) handleError(err error) error {
	return s.handleErrorWithFallback(err, bytecode.NoopOperand, nil)
}

// handleErrorWithFallback applies catch-table then protected-frame error policy
// and allows a catch-specific fallback assignment/action.
func (s *execState) handleErrorWithFallback(err error, dst bytecode.Operand, fallback runtime.Value) error {
	if err == nil {
		return nil
	}

	if catch, ok := s.tryCatch(s.pc); ok {
		if fallback != nil {
			s.registers.Values[dst] = fallback
		}

		if catch[2] >= 0 {
			s.pc = catch[2]
		}

		return nil
	}

	return s.handleProtectedError(err)
}

func (s *execState) wrapRuntimeError(err error) error {
	return diagnostic.WrapRuntimeError(s.program, s.pc, err)
}

func (s *execState) runtimeErrorFromPanic(r any) error {
	return diagnostic.RuntimeErrorFromPanic(s.program, s.pc, r)
}

func (s *execState) checkDivisionByZero(ctx context.Context, left, right runtime.Value) error {
	return diagnostic.CheckDivisionByZero(ctx, s.program, s.pc, left, right)
}

func (s *execState) checkModuloByZero(ctx context.Context, right runtime.Value) error {
	return diagnostic.CheckModuloByZero(ctx, s.program, s.pc, right)
}

func (s *execState) tryCatch(pos int) (bytecode.Catch, bool) {
	if s.catchByPC != nil && pos >= 0 && pos < len(s.catchByPC) {
		if idx := s.catchByPC[pos]; idx >= 0 {
			return s.program.CatchTable[idx], true
		}

		return bytecode.Catch{}, false
	}

	for _, pair := range s.program.CatchTable {
		if pos >= pair[0] && pos <= pair[1] {
			return pair, true
		}
	}

	return bytecode.Catch{}, false
}

func (s *execState) setOrTryCatch(dst bytecode.Operand, val runtime.Value, err error) error {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = val

		return nil
	}

	return s.handleErrorWithFallback(err, dst, runtime.None)
}

func (s *execState) setOrOptional(dst bytecode.Operand, val runtime.Value, err error, optional bool) error {
	if err == nil {
		s.registers.Values[dst] = val

		return nil
	}

	if optional || errors.Is(err, runtime.ErrNotFound) {
		s.registers.Values[dst] = runtime.None

		return nil
	}

	return err
}

func (s *execState) unwindToProtected() bool {
	registers, pc, ok := s.frames.UnwindToProtectedFrame(s.registers.Values)
	if !ok {
		return false
	}

	s.registers.Values = registers
	s.pc = pc
	return true
}

func (s *execState) returnToCaller(retVal runtime.Value) bool {
	registers, pc, ok := s.frames.ReturnToCaller(s.registers.Values, retVal)
	if !ok {
		return false
	}

	s.registers.Values = registers
	s.pc = pc
	return true
}
