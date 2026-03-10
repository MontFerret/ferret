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

type pendingFailure struct {
	err         error
	kind        errorKind
	mode        recoveryMode
	dst         bytecode.Operand
	fallback    runtime.Value
	setFallback bool
}

type execState struct {
	program     *bytecode.Program
	env         *Environment
	registers   *mem.RegisterFile
	scratch     *mem.Scratch
	frames      frame.CallStack
	catchByPC   []int
	pc          int
	panicPolicy PanicPolicy
	failure     pendingFailure
	hasFail     bool
}

func (s *execState) init(program *bytecode.Program, catchByPC []int, panicPolicy PanicPolicy) {
	s.program = program
	s.catchByPC = catchByPC
	s.panicPolicy = panicPolicy
	s.registers = mem.NewRegisterFile(program.Registers)
	s.scratch = mem.NewScratch(len(program.Params))
	s.frames.Init(maxUDFRegisters(program.Functions.UserDefined))
}

func (s *execState) prepareRun(env *Environment) {
	s.registers.Values = s.frames.Reset(s.registers.Values)

	if s.registers.IsDirty() {
		s.registers.Reset()
	}

	s.registers.MarkDirty()
	s.env = env
	s.pc = 0
	s.clearFailure()
}

func (s *execState) cleanupForPool() {
	s.registers.Values = s.frames.Reset(s.registers.Values)

	if s.registers.IsDirty() {
		s.registers.Reset()
	}

	for i := range s.scratch.Params {
		s.scratch.Params[i] = runtime.None
	}
	for i := range s.scratch.HostArgs {
		s.scratch.HostArgs[i] = runtime.None
	}

	s.env = nil
	s.pc = 0
	s.clearFailure()
}

func (s *execState) bindParams(env *Environment) error {
	required := s.program.Params

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

func (s *execState) raiseRuntime(err error, mode recoveryMode, dst bytecode.Operand, fallback runtime.Value, setFallback bool) {
	s.raise(err, errKindRuntime, mode, dst, fallback, setFallback)
}

func (s *execState) raiseInvariant(err error) {
	s.raise(err, errKindInvariant, recoverDefault, bytecode.NoopOperand, nil, false)
}

func (s *execState) raise(err error, kind errorKind, mode recoveryMode, dst bytecode.Operand, fallback runtime.Value, setFallback bool) {
	if err == nil || s.hasFail {
		return
	}

	s.failure = pendingFailure{
		err:         err,
		kind:        kind,
		mode:        mode,
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
		if s.panicPolicy == PanicPropagate {
			panic(failure.err)
		}

		return errReturn
	case errKindRuntime:
		switch failure.mode {
		case recoverOptional, recoverProtected:
			if failure.setFallback && failure.dst.IsRegister() {
				s.registers.Values[failure.dst] = normalizeValue(failure.fallback)
			}
			s.clearFailure()
			return errContinue
		default:
			if catch, ok := s.tryCatch(s.pc); ok {
				if failure.setFallback && failure.dst.IsRegister() {
					s.registers.Values[failure.dst] = normalizeValue(failure.fallback)
				}

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
	default:
		return errReturn
	}
}

func (s *execState) wrapRuntimeError(err error) error {
	return diagnostic.WrapRuntimeError(s.program, s.pc, s.frames.TraceEntries(), err)
}

func (s *execState) runtimeErrorFromPanic(r any) error {
	return diagnostic.RuntimeErrorFromPanic(s.program, s.pc, s.frames.TraceEntries(), r)
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

func (s *execState) setOrTryCatch(dst bytecode.Operand, val runtime.Value, err error) {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = normalizeValue(val)
		return
	}

	s.raiseRuntime(err, recoverDefault, dst, runtime.None, true)
}

func (s *execState) setOrOptional(dst bytecode.Operand, val runtime.Value, err error, optional bool) {
	if err == nil {
		s.registers.Values[dst] = normalizeValue(val)
		return
	}

	if optional || errors.Is(err, runtime.ErrNotFound) {
		s.raiseRuntime(err, recoverOptional, dst, runtime.None, true)
		return
	}

	s.raiseRuntime(err, recoverDefault, dst, runtime.None, true)
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

func (s *execState) setCallResult(op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = normalizeValue(out)
		return
	}

	if bytecode.IsProtectedCall(op) {
		s.raiseRuntime(err, recoverProtected, dst, runtime.None, true)
		return
	}

	s.raiseRuntime(err, recoverDefault, dst, runtime.None, true)
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
