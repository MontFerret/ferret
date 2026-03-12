package vm

import (
	"context"
	"errors"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type (
	pendingFailure struct {
		err         error
		fallback    runtime.Value
		pc          int
		dst         bytecode.Operand
		kind        errorKind
		mode        recoveryMode
		setFallback bool
	}

	execState struct {
		program   *bytecode.Program
		env       *Environment
		scratch   mem.Scratch
		frames    frame.CallStack
		catchByPC []int
		registers mem.RegisterFile
		failure   pendingFailure
		pc        int
		lastPC    int
		hasFail   bool
	}
)

func (s *execState) init(program *bytecode.Program, catchByPC []int) {
	s.program = program
	s.catchByPC = catchByPC
	s.registers.Init(program.Registers)
	s.scratch.Init(len(program.Params))
	s.frames.Init(maxUDFRegisters(program.Functions.UserDefined))
}

func (s *execState) start(env *Environment) {
	s.env = env
	s.pc = 0
	s.lastPC = -1
	s.clearFailure()
}

func (s *execState) end() {
	s.frames.Reset(s.registers.Values)
	s.registers.Reset()

	s.env = nil
	s.pc = 0
	s.lastPC = -1
	s.clearFailure()
}

func (s *execState) bindParams(env *Environment) error {
	required := s.program.Params
	if len(required) == 0 {
		return nil
	}

	s.scratch.ResizeParams(len(required))

	var missedParams []string

	for idx, name := range required {
		val, exists := env.Params[name]

		if !exists {
			if missedParams == nil {
				missedParams = make([]string, 0, len(required))
			}

			missedParams = append(missedParams, "@"+name)
			val = runtime.None
		}

		s.scratch.Params[idx] = val
	}

	if len(missedParams) > 0 {
		return runtime.Error(ErrMissedParam, strings.Join(missedParams, ", "))
	}

	return nil
}

func (s *execState) raiseRuntimeAt(pc int, err error, mode recoveryMode, dst bytecode.Operand, fallback runtime.Value, setFallback bool) {
	s.raise(pc, err, errKindRuntime, mode, dst, fallback, setFallback)
}

func (s *execState) raiseRuntime(err error, mode recoveryMode, dst bytecode.Operand, fallback runtime.Value, setFallback bool) {
	s.raiseRuntimeAt(s.pc, err, mode, dst, fallback, setFallback)
}

func (s *execState) raiseInvariantAt(pc int, err error) {
	s.raise(pc, err, errKindInvariant, recoverDefault, bytecode.NoopOperand, nil, false)
}

func (s *execState) raiseInvariant(err error) {
	s.raiseInvariantAt(s.pc, err)
}

func (s *execState) raise(pc int, err error, kind errorKind, mode recoveryMode, dst bytecode.Operand, fallback runtime.Value, setFallback bool) {
	if err == nil || s.hasFail {
		return
	}

	s.lastPC = pc
	s.failure = pendingFailure{
		err:         err,
		kind:        kind,
		mode:        mode,
		pc:          pc,
		dst:         dst,
		fallback:    fallback,
		setFallback: setFallback,
	}
	s.hasFail = true
}

func (s *execState) hasFailure() bool {
	return s.hasFail
}

func (s *execState) failureError() error {
	if !s.hasFail {
		return nil
	}

	return s.failure.err
}

func (s *execState) clearFailure() {
	if !s.hasFail {
		return
	}

	s.failure = pendingFailure{}
	s.hasFail = false
}

func (s *execState) resolveFailure() errAction {
	if !s.hasFail {
		return errOK
	}

	failure := s.failure

	switch failure.kind {
	case errKindInvariant:
		return errReturn
	case errKindRuntime:
		switch failure.mode {
		case recoverOptional:
			if s.isOptionalMiss(failure.err) {
				s.applyFailureFallback(failure)
				s.clearFailure()
				return errContinue
			}

			return s.resolveRuntimeDefault(failure)
		case recoverMissingMember:
			if s.isMissingMember(failure.err) {
				s.applyFailureFallback(failure)
				s.clearFailure()
				return errContinue
			}

			return s.resolveRuntimeDefault(failure)
		case recoverProtected:
			s.applyFailureFallback(failure)
			s.clearFailure()
			return errContinue
		default:
			return s.resolveRuntimeDefault(failure)
		}
	default:
		return errReturn
	}
}

func (s *execState) resolveRuntimeDefault(failure pendingFailure) errAction {
	if catch, ok := s.tryCatch(failure.pc); ok {
		s.applyFailureFallback(failure)

		if catch[2] >= 0 {
			s.pc = catch[2]
		}

		s.clearFailure()
		return errContinue
	}

	if s.unwindToProtected() {
		s.clearFailure()
		return errContinue
	}

	return errReturn
}

func (s *execState) applyFailureFallback(failure pendingFailure) {
	if failure.setFallback && failure.dst.IsRegister() {
		s.registers.Values[failure.dst] = normalizeValue(failure.fallback)
	}
}

func (s *execState) isOptionalMiss(err error) bool {
	return s.isMissingMember(err) || s.isNullMemberDereference(err)
}

func (s *execState) isMissingMember(err error) bool {
	return errors.Is(err, runtime.ErrNotFound)
}

func (s *execState) isNullMemberDereference(err error) bool {
	var memberErr *diagnostic.MemberAccessError
	if !errors.As(err, &memberErr) {
		return false
	}

	return runtime.IsSameType(memberErr.Target, runtime.TypeNone)
}

func (s *execState) errorPC() int {
	if s.hasFail {
		return s.failure.pc
	}

	if s.lastPC >= 0 {
		return s.lastPC
	}

	if s.pc > 0 {
		return s.pc - 1
	}

	return 0
}

func (s *execState) toRuntimePC(pc int) int {
	return pc + 1
}

func (s *execState) wrapRuntimeError(err error) error {
	return diagnostic.WrapRuntimeError(s.program, s.toRuntimePC(s.errorPC()), s.frames.TraceEntries(), err)
}

func (s *execState) runtimeErrorFromPanic(r any) error {
	return diagnostic.RuntimeErrorFromPanic(s.program, s.toRuntimePC(s.errorPC()), s.frames.TraceEntries(), r)
}

func (s *execState) checkDivisionByZeroAt(ctx context.Context, pc int, left, right runtime.Value) error {
	return diagnostic.CheckDivisionByZero(ctx, s.program, s.toRuntimePC(pc), left, right)
}

func (s *execState) checkModuloByZeroAt(ctx context.Context, pc int, right runtime.Value) error {
	return diagnostic.CheckModuloByZero(ctx, s.program, s.toRuntimePC(pc), right)
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

func (s *execState) setOrRaiseDefault(pc int, dst bytecode.Operand, val runtime.Value, err error) {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = normalizeValue(val)
		return
	}

	s.raiseRuntimeAt(pc, err, recoverDefault, dst, runtime.None, true)
}

func (s *execState) setOrOptional(pc int, dst bytecode.Operand, val runtime.Value, err error, optional bool) {
	if err == nil {
		s.registers.Values[dst] = normalizeValue(val)
		return
	}

	if optional && s.isNullMemberDereference(err) {
		s.registers.Values[dst] = runtime.None
		return
	}

	mode := recoverMissingMember
	if optional {
		mode = recoverOptional
	}

	s.raiseRuntimeAt(pc, err, mode, dst, runtime.None, true)
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

func (s *execState) setCallResult(pc int, op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = normalizeValue(out)
		return
	}

	if bytecode.IsProtectedCall(op) {
		s.raiseRuntimeAt(pc, err, recoverProtected, dst, runtime.None, true)
		return
	}

	s.raiseRuntimeAt(pc, err, recoverDefault, dst, runtime.None, true)
}

func (s *execState) resolveUdfID(val runtime.Value) (int, error) {
	idVal, ok := val.(runtime.Int)
	if !ok {
		return -1, ErrInvalidFunctionName
	}

	return int(idVal), nil
}

func (s *execState) udfByID(id int) (*bytecode.UDF, error) {
	if id < 0 || s.program == nil || id >= len(s.program.Functions.UserDefined) {
		return nil, ErrUnresolvedFunction
	}

	return &s.program.Functions.UserDefined[id], nil
}
