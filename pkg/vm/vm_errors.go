package vm

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
)

type (
	errAction uint8
	errPolicy uint8
)

const (
	errOK errAction = iota
	errContinue
	errReturn
)

const (
	errPolicyProtected errPolicy = iota
	errPolicyCatch
)

func (s *execState) applyErrPolicy(err error, policy errPolicy, dst bytecode.Operand, fallback runtime.Value) errAction {
	if err == nil {
		return errOK
	}

	return s.applyErrPolicySlow(err, policy, dst, fallback)
}

func (s *execState) applyErrPolicySlow(err error, policy errPolicy, dst bytecode.Operand, fallback runtime.Value) errAction {
	switch policy {
	case errPolicyCatch:
		if catch, ok := s.tryCatch(s.pc); ok {
			if fallback != nil {
				s.registers.Values[dst] = fallback
			}

			if catch[2] >= 0 {
				s.pc = catch[2]
			}

			return errContinue
		}
		fallthrough
	case errPolicyProtected:
		if s.unwindToProtected() {
			return errContinue
		}

		return errReturn
	default:
		if s.unwindToProtected() {
			return errContinue
		}

		return errReturn
	}
}

func (s *execState) applyProtected(err error) errAction {
	return s.applyErrPolicy(err, errPolicyProtected, bytecode.NoopOperand, nil)
}

func (s *execState) applyCatch(dst bytecode.Operand, fallback runtime.Value, err error) errAction {
	return s.applyErrPolicy(err, errPolicyCatch, dst, fallback)
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

func (s *execState) setOrTryCatch(dst bytecode.Operand, val runtime.Value, err error) errAction {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = val
		return errOK
	}

	return s.applyErrPolicySlow(err, errPolicyCatch, dst, runtime.None)
}

func (s *execState) setOrOptional(dst bytecode.Operand, val runtime.Value, err error, optional bool) errAction {
	if err == nil {
		s.registers.Values[dst] = val
		return errOK
	}

	if optional || errors.Is(err, runtime.ErrNotFound) {
		s.registers.Values[dst] = runtime.None
		return errContinue
	}

	return errReturn
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
