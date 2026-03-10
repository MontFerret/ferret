package vm

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type execState struct {
	program     *bytecode.Program
	env         *Environment
	registers   *mem.RegisterFile
	scratch     *mem.Scratch
	frames      frame.CallStack
	catchByPC   []int
	pc          int
	panicPolicy PanicPolicy
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
}

func (s *execState) cleanupForPool() {
	s.registers.Values = s.frames.Reset(s.registers.Values)

	if s.registers.IsDirty() {
		s.registers.Reset()
	}

	for i := range s.scratch.Params {
		s.scratch.Params[i] = runtime.None
	}

	s.env = nil
	s.pc = 0
}

func (s *execState) fail(err error, class failClass, dst bytecode.Operand, fallback runtime.Value, setFallback bool) errAction {
	if err == nil {
		return errOK
	}

	switch class {
	case failInvariant:
		if s.panicPolicy == PanicPropagate {
			panic(err)
		}

		return errReturn
	case failProtected:
		if setFallback && dst.IsRegister() {
			s.registers.Values[dst] = normalizeValue(fallback)
		}

		return errContinue
	default:
		if catch, ok := s.tryCatch(s.pc); ok {
			if setFallback && dst.IsRegister() && fallback != nil {
				s.registers.Values[dst] = normalizeValue(fallback)
			}

			if catch[2] >= 0 {
				s.pc = catch[2]
			}

			return errContinue
		}

		if s.unwindToProtected() {
			return errContinue
		}

		return errReturn
	}
}

func (s *execState) applyProtected(err error) errAction {
	return s.fail(err, failRuntime, bytecode.NoopOperand, nil, false)
}

func (s *execState) applyCatch(dst bytecode.Operand, fallback runtime.Value, err error) errAction {
	return s.fail(err, failRuntime, dst, fallback, true)
}

func (s *execState) applyInvariant(err error) errAction {
	return s.fail(err, failInvariant, bytecode.NoopOperand, nil, false)
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

func (s *execState) setOrTryCatch(dst bytecode.Operand, val runtime.Value, err error) errAction {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = normalizeValue(val)
		return errOK
	}

	return s.fail(err, failRuntime, dst, runtime.None, true)
}

func (s *execState) setOrOptional(dst bytecode.Operand, val runtime.Value, err error, optional bool) errAction {
	if err == nil {
		s.registers.Values[dst] = normalizeValue(val)
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

func (s *execState) setCallResult(op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) errAction {
	reg := s.registers.Values

	if err == nil {
		reg[dst] = normalizeValue(out)

		return errOK
	}

	if bytecode.IsProtectedCall(op) {
		return s.fail(err, failProtected, dst, runtime.None, true)
	}

	return s.fail(err, failRuntime, dst, runtime.None, true)
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
