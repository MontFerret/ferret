package vm

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
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
		windows   mem.WindowPool
		catchByPC []int
		registers mem.RegisterFile
		owned     mem.OwnedResources
		aliases   mem.AliasTracker
		deferred  mem.DeferredClosers
		failure   pendingFailure
		pc        int
		lastPC    int
		hasFail   bool
	}
)

func (s *execState) init(program *bytecode.Program) {
	s.program = program
	s.catchByPC = buildCatchByPC(len(program.Bytecode), program.CatchTable)
	s.registers = mem.NewRegisterFile(program.Registers)
	s.scratch = mem.NewScratch(len(program.Params))
	s.frames = frame.NewCallStack()
	s.windows = mem.NewWindowPool(maxUDFRegisters(program.Functions.UserDefined))
}

func (s *execState) startRun(env *Environment) error {
	s.env = env
	s.owned.Reset()
	s.aliases.Reset()
	s.deferred.Reset()
	s.pc = 0
	s.lastPC = -1
	s.clearFailure()

	return s.bindParams(env)
}

func (s *execState) endRun() {
	for {
		s.owned.DrainTo(&s.deferred)

		caller, ok := s.frames.Pop()
		if !ok {
			break
		}

		s.windows.Release(s.registers)
		s.registers = caller.CallerRegisters
		s.owned = caller.OwnedResources
		s.aliases = caller.Aliases
	}

	_ = s.deferred.CloseAll()
	s.resetRunStorage()
}

func (s *execState) finishRun(root runtime.Value) *Result {
	return s.finishRunInto(root, newResult(root))
}

func (s *execState) finishRunInto(root runtime.Value, result *Result) *Result {
	result.reset(root)

	if s.owned.Empty() && s.deferred.Empty() {
		s.resetRunStorage()
		return result
	}

	var resultOwned mem.OwnedResources

	if key, closer, ok := mem.ResourceKeyOf(root); ok && s.owned.ExtractByKey(key) {
		resultOwned.TrackResolved(key, closer)
	}
	if !s.owned.Empty() {
		s.owned.DrainTo(&s.deferred)
	}

	if !resultOwned.Empty() {
		result.adoptOwned(&resultOwned)
	}
	if !s.deferred.Empty() {
		result.adoptDeferred(&s.deferred)
	}

	s.resetRunStorage()

	return result
}

func (s *execState) resetRunStorage() {
	s.owned.Reset()
	s.aliases.Reset()
	s.deferred.Reset()

	s.registers.Reset()
	s.scratch.Reset()

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

		val = normalizeValue(val)
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
		s.writeBorrowedRegister(failure.dst, failure.fallback)
	}
}

func (s *execState) isOptionalMiss(err error) bool {
	return s.isMissingMember(err) || s.isNullMemberDereference(err)
}

func (s *execState) isMissingMember(err error) bool {
	return errors.Is(err, runtime.ErrNotFound)
}

func (s *execState) isNullMemberDereference(err error) bool {
	var memberErr *diagnostics.MemberAccessError
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
	return diagnostics.WrapRuntimeError(s.program, s.toRuntimePC(s.errorPC()), s.frames.TraceEntries(), err)
}

func (s *execState) runtimeErrorFromPanic(r any) error {
	return diagnostics.RuntimeErrorFromPanic(s.program, s.toRuntimePC(s.errorPC()), s.frames.TraceEntries(), r)
}

func (s *execState) checkDivisionByZeroAt(ctx context.Context, pc int, left, right runtime.Value) error {
	return diagnostics.CheckDivisionByZero(ctx, s.program, s.toRuntimePC(pc), left, right)
}

func (s *execState) checkModuloByZeroAt(ctx context.Context, pc int, right runtime.Value) error {
	return diagnostics.CheckModuloByZero(ctx, s.program, s.toRuntimePC(pc), right)
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
	if err == nil {
		s.writeBorrowedRegister(dst, val)
		return
	}

	s.raiseRuntimeAt(pc, err, recoverDefault, dst, runtime.None, true)
}

func (s *execState) setProducedOrRaiseDefault(pc int, dst bytecode.Operand, val runtime.Value, err error) {
	if err == nil {
		s.writeProducedRegister(dst, val)
		return
	}

	s.raiseRuntimeAt(pc, err, recoverDefault, dst, runtime.None, true)
}

func (s *execState) setOrOptional(pc int, dst bytecode.Operand, val runtime.Value, err error, optional bool) {
	if err == nil {
		s.writeBorrowedRegister(dst, val)
		return
	}

	if optional && s.isNullMemberDereference(err) {
		s.writeBorrowedRegister(dst, runtime.None)
		return
	}

	mode := recoverMissingMember
	if optional {
		mode = recoverOptional
	}

	s.raiseRuntimeAt(pc, err, mode, dst, runtime.None, true)
}

func (s *execState) unwindToProtected() bool {
	boundary := s.frames.NearestRecoveryBoundary()
	if boundary < 0 {
		return false
	}

	for s.frames.Len() > boundary+1 {
		frame, ok := s.frames.Pop()
		if !ok {
			return false
		}

		s.owned.DrainTo(&s.deferred)
		s.windows.Release(s.registers)
		s.registers = frame.CallerRegisters
		s.owned = frame.OwnedResources
		s.aliases = frame.Aliases
	}

	frame, ok := s.frames.Pop()
	if !ok {
		return false
	}

	s.owned.DrainTo(&s.deferred)
	s.windows.Release(s.registers)
	s.registers = frame.CallerRegisters
	s.owned = frame.OwnedResources
	s.aliases = frame.Aliases
	if frame.ReturnDest.IsRegister() {
		s.writeBorrowedRegister(frame.ReturnDest, runtime.None)
	}
	s.pc = frame.ReturnPC
	return true
}

func (s *execState) returnToCaller(retVal runtime.Value) bool {
	frame, ok := s.frames.Pop()
	if !ok {
		return false
	}

	var (
		retKey    mem.ResourceKey
		retCloser io.Closer
		retOwned  bool
	)

	if !s.owned.Empty() {
		if key, closer, ok := mem.ResourceKeyOf(retVal); ok {
			retKey = key
			retCloser = closer
			retOwned = s.owned.ExtractByKey(key)
		}

		s.owned.DrainTo(&s.deferred)
	}

	s.windows.Release(s.registers)
	s.registers = frame.CallerRegisters
	s.owned = frame.OwnedResources
	s.aliases = frame.Aliases
	if frame.ReturnDest.IsRegister() {
		s.writeBorrowedRegister(frame.ReturnDest, retVal)

		if retOwned {
			s.owned.TrackResolved(retKey, retCloser)
			s.aliases.Inc(retKey)
		}
	}
	s.pc = frame.ReturnPC
	return true
}

func (s *execState) setCallResult(pc int, op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) {
	if err == nil {
		s.writeProducedRegister(dst, out)
		return
	}

	if bytecode.IsProtectedCall(op) {
		s.raiseRuntimeAt(pc, err, recoverProtected, dst, runtime.None, true)
		return
	}

	s.raiseRuntimeAt(pc, err, recoverDefault, dst, runtime.None, true)
}

// decAliasAndMaybeDiscard decrements the alias count for key and, if no
// live register aliases remain, discards the resource from owned resources.
// This replaces the old O(n) register scan with an O(1) map lookup.
func (s *execState) decAliasAndMaybeDiscard(key mem.ResourceKey, closer io.Closer) {
	if s.aliases.Dec(key) > 0 {
		return
	}

	s.owned.DiscardByKey(key, closer, &s.deferred)
}

// retireOwnership terminally removes VM ownership of val after the value
// escapes into an external sink. Any stale register aliases are ignored by
// later automatic cleanup.
func (s *execState) retireOwnership(val runtime.Value) {
	key, _, ok := mem.ResourceKeyOf(val)
	if !ok {
		return
	}

	s.owned.ExtractByKey(key)
	s.aliases.Delete(key)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Register Write Discipline
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
//
// These are the ONLY legal ways to mutate register state.
// Direct reg[i] = v assignments are forbidden outside these helpers.
//
// Each helper encodes ownership semantics and ensures correct cleanup:
//
//   writeBorrowedRegister(dst, val)   - Write a borrowed/untracked value
//   writeProducedRegister(dst, val)   - Write a produced/owned value
//   copyRegister(dst, src)            - Copy register value (ownership follows value)
//   clearRegister(dst)                - Clear register to None
//
// Hot path optimization:
//   - Most register writes don't involve closers
//   - Check if old value is owned FIRST (cheap map lookup)
//   - Only scan for aliases when overwriting an owned closer
//   - Early return on common non-closer case
//
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

// writeBorrowedRegister writes a borrowed (not owned) value to a register.
// Use this for values that are:
//   - constants
//   - parameters
//   - computed values that don't need tracking (numbers, strings, bools)
//   - values from other registers (borrowed references)
//
// Fast path: obvious scalar values cannot participate in ownership, so we can
// overwrite them without probing ownership state.
func (s *execState) writeBorrowedRegister(dst bytecode.Operand, val runtime.Value) runtime.Value {
	val = normalizeValue(val)

	if !dst.IsRegister() {
		return val
	}

	prev := s.registers[dst]

	if !mem.CanTrackValue(prev) {
		s.registers[dst] = val
		return val
	}

	prevKey, prevCloser, ok := mem.ResourceKeyOf(prev)
	if !ok || !s.owned.OwnsKey(prevKey) {
		s.registers[dst] = val
		return val
	}

	newKey, _, newOK := mem.ResourceKeyOf(val)

	if !newOK || prevKey != newKey {
		s.decAliasAndMaybeDiscard(prevKey, prevCloser)
	}

	s.registers[dst] = val
	return val
}

// writeProducedRegister writes a produced (owned) value to a register.
// Use this for values that are:
//   - newly allocated objects (arrays, objects, iterators)
//   - returned from host functions
//   - created by VM instructions (streams, collectors, etc.)
//
// The value will be tracked for ownership and cleanup when it resolves to a
// closable resource.
func (s *execState) writeProducedRegister(dst bytecode.Operand, val runtime.Value) runtime.Value {
	val = normalizeValue(val)
	newKey, newCloser, newOK := mem.ResourceKeyOf(val)

	if !dst.IsRegister() {
		if newOK {
			s.owned.TrackResolved(newKey, newCloser)
		}

		return val
	}

	prev := s.registers[dst]

	if !mem.CanTrackValue(prev) {
		s.registers[dst] = val
		if newOK {
			s.owned.TrackResolved(newKey, newCloser)
			s.aliases.Inc(newKey)
		}

		return val
	}

	prevKey, prevCloser, ok := mem.ResourceKeyOf(prev)
	if !ok || !s.owned.OwnsKey(prevKey) {
		s.registers[dst] = val
		if newOK {
			s.owned.TrackResolved(newKey, newCloser)
			s.aliases.Inc(newKey)
		}

		return val
	}

	if !newOK || prevKey != newKey {
		s.decAliasAndMaybeDiscard(prevKey, prevCloser)
	}

	s.registers[dst] = val

	if newOK {
		s.owned.TrackResolved(newKey, newCloser)
		s.aliases.Inc(newKey)
	}

	return val
}

// copyRegister moves a value from one register to another, transferring ownership if any.
// Use this for explicit register moves (OpMove) where:
//   - source register will be overwritten or cleared soon
//   - you want to avoid aliasing complexity
//
// This is semantically a copy, but ownership stays with the value, not the slot.
// If the copied value is an owned closer, its alias count is incremented.
func (s *execState) copyRegister(dst, src bytecode.Operand) runtime.Value {
	if !src.IsRegister() {
		return s.writeBorrowedRegister(dst, runtime.None)
	}

	if !dst.IsRegister() {
		return s.registers[src]
	}

	val := s.registers[src]
	prev := s.registers[dst]
	valKey, _, valTrackable := mem.ResourceKeyOf(val)
	valOwned := valTrackable && s.owned.OwnsKey(valKey)

	if !mem.CanTrackValue(prev) {
		s.registers[dst] = val
		if valOwned {
			s.aliases.Inc(valKey)
		}

		return val
	}

	prevKey, prevCloser, ok := mem.ResourceKeyOf(prev)
	if !ok || !s.owned.OwnsKey(prevKey) {
		s.registers[dst] = val
		if valOwned {
			s.aliases.Inc(valKey)
		}

		return val
	}

	if !valOwned || prevKey != valKey {
		s.decAliasAndMaybeDiscard(prevKey, prevCloser)
	}

	s.registers[dst] = val
	if valOwned {
		s.aliases.Inc(valKey)
	}

	return val
}

// clearRegister clears a register, setting it to None and discarding any owned value.
// Use this for explicit register cleanup when you don't need to set a new value.
func (s *execState) clearRegister(dst bytecode.Operand) {
	if !dst.IsRegister() {
		return
	}

	prev := s.registers[dst]

	if !mem.CanTrackValue(prev) {
		s.registers[dst] = runtime.None
		return
	}

	if key, closer, ok := mem.ResourceKeyOf(prev); ok && s.owned.OwnsKey(key) {
		s.decAliasAndMaybeDiscard(key, closer)
	}

	s.registers[dst] = runtime.None
}

func (s *execState) udfByID(id int) (*bytecode.UDF, error) {
	if id < 0 || s.program == nil || id >= len(s.program.Functions.UserDefined) {
		return nil, ErrUnresolvedFunction
	}

	return &s.program.Functions.UserDefined[id], nil
}
