package debugger

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"

	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Session controls one retained-state source-level debug execution.
// Execution and inspection commands are serialized. Breakpoint mutation,
// breakpoint listing, Pause, and Close are safe during a running command.
type Session struct {
	values           vm.DebugValueAccess
	services         SessionServices
	closeErr         error
	runCtx           context.Context
	executionCtx     context.Context
	execution        vm.DebugExecution
	valueRefs        map[ValueReference]runtime.Value
	runCancel        context.CancelCauseFunc
	breakpointCheck  vm.DebugBreakpointPredicate
	breakpointWrite  chan struct{}
	breakpointDone   chan struct{}
	breakpointState  string         // lifecycleMu; empty until native terminal commitment
	hitBreakpointIDs []BreakpointID // commandMu; captured by the execution callback
	breakpointData   atomic.Pointer[breakpointSnapshot]
	activeCancel     context.CancelCauseFunc
	params           []string
	pointIndex       debugpoint.Index
	source           source.Source
	format           FormatOptions
	nextValueRef     ValueReference
	activeCancelID   uint64
	closeOnce        sync.Once
	commandMu        sync.Mutex
	lifecycleMu      sync.Mutex
	closed           atomic.Bool
	started          atomic.Bool
	afterRun         bool
	executionClosed  bool
}

// NewSession creates an advanced debugger session from explicit dependencies.
func NewSession(config Config) (*Session, error) {
	if config.Execution == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debug execution is required")
	}

	if config.Values == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debug value access is required")
	}

	if config.Services == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debug session services are required")
	}

	if config.Source.Empty() {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debug source is required")
	}

	format := config.Format
	if format.MaxDepth <= 0 || format.MaxItems <= 0 || format.MaxBytes <= 0 {
		format = DefaultFormatOptions()
	}

	pointIndex, err := debugpoint.New(config.DebugPoints)
	if err != nil {
		return nil, runtime.Errorf(runtime.ErrInvalidArgument, "invalid debug points: %s", err)
	}

	session := &Session{
		execution:       config.Execution,
		values:          config.Values,
		services:        config.Services,
		source:          config.Source,
		pointIndex:      pointIndex,
		params:          append([]string(nil), config.Params...),
		format:          format,
		breakpointWrite: make(chan struct{}, 1),
		breakpointDone:  make(chan struct{}),
		valueRefs:       make(map[ValueReference]runtime.Value),
		nextValueRef:    1,
	}
	session.breakpointData.Store(&breakpointSnapshot{nextID: 1})
	session.breakpointCheck = session.hasBreakpoint

	return session, nil
}

// Start begins execution and stops at the first executable source location.
func (s *Session) Start(ctx context.Context) (*Event, error) {
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}

	if s.started.Load() {
		return nil, &StateError{Operation: "start", State: "started"}
	}

	runCtx, err := s.services.BeforeRun(ctx)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("before run hooks: %w", err), ctx.Err())
	}

	err = s.checkContext(runCtx)
	if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
		err = errors.Join(err, ctxErr)
	}

	if err != nil {
		if runCtx == nil {
			runCtx = ctx
		}

		// Settle this attempt without consuming a later successful Start's hooks
		// or making Close repeat them. Execution has not been entered yet.
		if hookErr := s.services.AfterRun(runCtx, err); hookErr != nil {
			err = errors.Join(err, fmt.Errorf("after run hooks: %w", hookErr))
		}

		return nil, err
	}

	executionCtx, runCancel := context.WithCancelCause(runCtx)

	s.started.Store(true)
	s.runCtx = runCtx
	s.executionCtx = executionCtx
	s.setRunCancel(runCancel)
	s.resetValueReferences()

	event, err := s.execution.Start(s.services.ExtendContext(executionCtx))
	if err != nil {
		s.recordExecutionTerminal()

		return nil, err
	}

	return s.convertEvent(event)
}

// Continue resumes execution until a breakpoint, pause request, error, or completion.
func (s *Session) Continue(ctx context.Context) (*Event, error) {
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	return s.resume(ctx, vm.DebugResumeContinue)
}

// StepIn resumes execution until the next debuggable source location and may
// enter called functions.
func (s *Session) StepIn(ctx context.Context) (*Event, error) {
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	return s.resume(ctx, vm.DebugResumeStepIn)
}

// StepOver resumes execution until the next debuggable source location at the
// same or shallower call depth. Breakpoints and pause requests may interrupt it
// inside deeper calls.
func (s *Session) StepOver(ctx context.Context) (*Event, error) {
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	return s.resume(ctx, vm.DebugResumeStepOver)
}

// StepOut resumes execution until the current frame exits, then stops at the
// next debuggable source location in a caller. At main it runs to completion.
// Breakpoints and pause requests may interrupt it before the frame exits.
func (s *Session) StepOut(ctx context.Context) (*Event, error) {
	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.checkContext(ctx); err != nil {
		return nil, err
	}

	return s.resume(ctx, vm.DebugResumeStepOut)
}

// Pause requests a stop at the next logical source location.
func (s *Session) Pause() error {
	if err := s.ensureOpen(); err != nil {
		return err
	}

	if !s.started.Load() {
		return &StateError{Operation: "pause", State: "new"}
	}
	s.execution.RequestPause()

	return nil
}

// Frames returns the current frame followed by callers.
func (s *Session) Frames() ([]Frame, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	frames, err := s.execution.Frames()

	if err != nil {
		return nil, err
	}
	out := make([]Frame, 0, len(frames))

	for _, frame := range frames {
		out = append(out, Frame{
			Name:       frame.Name,
			FunctionID: apidebugger.FunctionID(frame.FunctionID),
			Location:   s.locationForPC(frame.PC, frame.FunctionID),
		})
	}

	return out, nil
}

// Locals returns the visible top-frame locals followed by bound parameters.
func (s *Session) Locals() ([]Variable, error) {
	return s.FrameLocals(0)
}

// FrameLocals returns the visible locals and bound parameters for one paused
// frame. Frame indexes follow Frames: zero is the current frame.
func (s *Session) FrameLocals(frame int) ([]Variable, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	locals, err := s.frameLocals(frame)

	if err != nil {
		return nil, err
	}
	out := make([]Variable, 0, len(locals)+len(s.execution.Params()))

	for _, local := range locals {
		out = append(out, Variable{Name: local.Name, Mutable: local.Mutable, Value: s.debugValue(local.Value)})
	}

	params := s.execution.Params()
	names := append([]string(nil), s.params...)

	sort.Strings(names)

	for _, name := range names {
		if value, exists := params.Get(name); exists {
			out = append(out, Variable{Name: "@" + name, Param: true, Value: s.debugValue(value)})
		}
	}

	return out, nil
}

func (s *Session) frameLocals(frame int) ([]vm.DebugLocal, error) {
	if inspector, ok := s.execution.(vm.DebugFrameInspector); ok {
		return inspector.FrameLocals(frame)
	}

	if frame != 0 {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "debug execution does not support caller frame inspection")
	}

	return s.execution.Locals()
}

// Variables returns the child variables for one expandable debugger value from
// the current paused state.
func (s *Session) Variables(reference ValueReference) ([]Variable, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}

	if !reference.Valid() {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debug value reference must be positive")
	}

	value, exists := s.valueRefs[reference]
	if !exists {
		return nil, runtime.Errorf(runtime.ErrNotFound, "debug value reference %d", reference)
	}

	inspection, ok := s.expandableInspection(value)
	if !ok {
		return nil, runtime.Errorf(runtime.ErrInvalidOperation, "debug value reference %d is not expandable", reference)
	}

	if inspection.Kind == vm.DebugValueArray {
		out := make([]Variable, 0, len(inspection.Items))

		for i, item := range inspection.Items {
			out = append(out, Variable{
				Name:  strconv.Itoa(i),
				Value: s.debugValue(item.Value),
			})
		}

		return out, nil
	}

	items := append([]vm.DebugValueItem(nil), inspection.Items...)
	sort.Slice(items, func(i, j int) bool { return items[i].Key < items[j].Key })

	out := make([]Variable, 0, len(items))
	for _, item := range items {
		out = append(out, Variable{
			Name:  item.Key,
			Value: s.debugValue(item.Value),
		})
	}

	return out, nil
}

// Evaluate evaluates a conservative, side-effect-free expression against the
// paused top frame. It supports literals, locals, parameters, supported member
// reads, scalar arithmetic/comparisons, boolean logic, and conditionals. It
// rejects calls, queries, mutation, async/event behavior, and full collection
// execution.
func (s *Session) Evaluate(ctx context.Context, expression string) (Value, error) {
	return s.EvaluateFrame(ctx, 0, expression)
}

// EvaluateFrame evaluates an expression against one paused frame.
func (s *Session) EvaluateFrame(ctx context.Context, frame int, expression string) (Value, error) {
	if err := s.checkContext(ctx); err != nil {
		return Value{}, err
	}

	if err := s.lockCommand(); err != nil {
		return Value{}, err
	}
	defer s.commandMu.Unlock()

	if err := s.checkContext(ctx); err != nil {
		return Value{}, err
	}

	if err := s.ensureOpen(); err != nil {
		return Value{}, err
	}
	locals, err := s.frameLocals(frame)

	if err != nil {
		return Value{}, err
	}
	values := make(map[string]runtime.Value, len(locals))

	for _, local := range locals {
		values[local.Name] = local.Value
	}

	value, err := evaluateExpression(ctx, expression, evalScope{
		locals: values,
		params: s.execution.Params(),
		values: s.values,
	})

	if err != nil {
		return Value{}, err
	}

	return s.debugValue(value), nil
}

func (s *Session) resume(ctx context.Context, mode vm.DebugResumeMode) (*Event, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}

	if !s.started.Load() {
		return nil, &StateError{Operation: "resume", State: "new"}
	}

	switch s.execution.Status() {
	case vm.DebugExecutionCompleted:
		return nil, &StateError{Operation: "resume", State: "completed"}
	case vm.DebugExecutionTerminated:
		return nil, &StateError{Operation: "resume", State: "terminated"}
	}

	if ctx == s.runCtx || ctx == s.executionCtx {
		ctx = s.executionCtx
	} else {
		var cleanup func()
		var cancel context.CancelCauseFunc
		ctx, cleanup, cancel = newResumeContext(s.runCtx, ctx)
		defer cleanup()
		activeID := s.setActiveCancel(cancel)
		defer s.clearActiveCancel(activeID)
	}

	s.resetValueReferences()

	s.hitBreakpointIDs = nil
	event, err := s.execution.Resume(s.services.ExtendContext(ctx), mode, s.breakpointCheck)
	if err != nil {
		s.recordExecutionTerminal()

		return nil, err
	}

	return s.convertEvent(event)
}

func (s *Session) convertEvent(event *vm.DebugExecutionEvent) (*Event, error) {
	if event == nil {
		return nil, runtime.Error(runtime.ErrUnexpected, "debug execution returned no event")
	}

	if event.Reason == vm.DebugStopCompleted {
		s.finishBreakpoints("completed")
	} else if event.Reason == vm.DebugStopTerminated {
		s.finishBreakpoints("terminated")
	}

	out := &Event{Depth: event.Depth, Error: event.Error}

	if event.Point != nil {
		out.Location = s.source.RangeAt(event.Point.Span)
	}

	if s.closed.Load() && event.Reason != vm.DebugStopCompleted && event.Reason != vm.DebugStopTerminated {
		out.Reason = ReasonTerminated
		out.Error = context.Cause(s.executionCtx)

		if out.Error == nil {
			out.Error = context.Canceled
		}

		if event.Error != nil {
			out.Error = errors.Join(out.Error, event.Error)
		}

		// Close can win while the VM is returning an inspectable stop. Drain
		// retained state before publishing termination and cache the cleanup
		// result without making the pending Close call the execution again.
		s.executionClosed = true
		if closeErr := s.execution.Close(); closeErr != nil {
			s.closeErr = errors.Join(s.closeErr, closeErr)
			out.Error = errors.Join(out.Error, closeErr)
		}

		if hookErr := s.runAfterHooks(out.Error); hookErr != nil {
			out.Error = errors.Join(out.Error, hookErr)
		}

		return out, nil
	}

	switch event.Reason {
	case vm.DebugStopEntry:
		out.Reason = ReasonEntry
	case vm.DebugStopBreakpoint:
		out.Reason = ReasonBreakpoint
		out.HitBreakpointIDs = append([]BreakpointID(nil), s.hitBreakpointIDs...)
	case vm.DebugStopStep:
		out.Reason = ReasonStep
	case vm.DebugStopPause:
		out.Reason = ReasonPause
	case vm.DebugStopRuntimeError:
		out.Reason = ReasonRuntimeError
		if hookErr := s.runAfterHooks(event.Error); hookErr != nil {
			out.Error = errors.Join(out.Error, hookErr)
		}
	case vm.DebugStopCompleted:
		s.resetValueReferences()
		out.Reason = ReasonCompleted
		output, outputErr := s.services.Materialize(event.Result)
		closeErr := event.Result.Close()
		hookErr := s.runAfterHooks(nil)

		out.Output = output

		if outputErr != nil || closeErr != nil || hookErr != nil {
			return out, errors.Join(outputErr, closeErr, hookErr)
		}
	case vm.DebugStopTerminated:
		s.resetValueReferences()
		out.Reason = ReasonTerminated

		if hookErr := s.runAfterHooks(event.Error); hookErr != nil {
			out.Error = errors.Join(out.Error, hookErr)
		}
	}

	return out, nil
}

func (s *Session) runAfterHooks(runErr error) error {
	if s.afterRun {
		return nil
	}

	s.afterRun = true

	ctx := s.runCtx
	if s.executionCtx != nil {
		ctx = s.executionCtx
	}

	if hookErr := s.services.AfterRun(ctx, runErr); hookErr != nil {
		return fmt.Errorf("after run hooks: %w", hookErr)
	}

	return nil
}

func (s *Session) debugValue(value runtime.Value) Value {
	info, _ := s.values.DebugInfo(value)
	typeName := info.TypeName

	if typeName == "" {
		typeName = s.values.TypeName(value)
		info.TypeName = typeName
	} else {
		typeName = boundedText(typeName, s.format.MaxBytes)
	}

	return Value{
		Reference: s.referenceForValue(value),
		Type:      typeName,
		Display:   formatValueWithInfo(value, info, s.values, s.format),
	}
}

func (s *Session) locationForPC(pc int, functionID bytecode.FunctionID) source.Location {
	var point *bytecode.DebugPoint

	if functionID == bytecode.NoFunction || functionID.Valid() {
		point = s.pointIndex.NearestBeforeOrAtInFunction(functionID, pc)
	} else {
		point = s.pointIndex.NearestBeforeOrAt(pc)
	}

	if point == nil {
		return source.Location{SourceName: s.source.Name()}
	}

	return s.source.LocationAt(point.Span)
}

func (s *Session) ensureOpen() error {
	if s == nil || s.closed.Load() {
		return &StateError{Operation: "use", State: "closed"}
	}

	return nil
}

func (s *Session) lockCommand() error {
	if s == nil || s.closed.Load() {
		return &StateError{Operation: "use", State: "closed"}
	}

	s.commandMu.Lock()
	if s.closed.Load() {
		s.commandMu.Unlock()
		return &StateError{Operation: "use", State: "closed"}
	}

	return nil
}

func (s *Session) setRunCancel(cancel context.CancelCauseFunc) {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()

	s.runCancel = cancel
	if s.closed.Load() && s.runCancel != nil {
		s.runCancel(context.Canceled)
	}
}

func (s *Session) setActiveCancel(cancel context.CancelCauseFunc) uint64 {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()

	s.activeCancelID++
	s.activeCancel = cancel
	if s.closed.Load() && s.activeCancel != nil {
		s.activeCancel(context.Canceled)
	}

	return s.activeCancelID
}

func (s *Session) clearActiveCancel(id uint64) {
	s.lifecycleMu.Lock()
	defer s.lifecycleMu.Unlock()

	if s.activeCancelID == id {
		s.activeCancel = nil
	}
}

func (s *Session) requestTermination() {
	s.lifecycleMu.Lock()

	if s.runCancel != nil {
		s.runCancel(context.Canceled)
	}

	if s.activeCancel != nil {
		s.activeCancel(context.Canceled)
	}

	s.lifecycleMu.Unlock()

	if s.execution != nil {
		s.execution.RequestPause()
	}
}

func (s *Session) resetValueReferences() {
	if len(s.valueRefs) == 0 {
		return
	}

	clear(s.valueRefs)
}

func (s *Session) referenceForValue(value runtime.Value) ValueReference {
	if _, ok := s.expandableInspection(value); !ok {
		return 0
	}

	reference := s.nextValueRef
	s.nextValueRef++
	s.valueRefs[reference] = value

	return reference
}

func (s *Session) expandableInspection(value runtime.Value) (vm.DebugValueInspection, bool) {
	inspection, ok := s.values.Inspect(value, s.format.MaxItems)
	if !ok || !inspection.Complete || inspection.Length <= 0 || inspection.Length > s.format.MaxItems {
		return vm.DebugValueInspection{}, false
	}

	if inspection.Kind != vm.DebugValueArray && inspection.Kind != vm.DebugValueObject {
		return vm.DebugValueInspection{}, false
	}

	return inspection, true
}

// Close requests termination, then releases retained VM state and
// embedding-owned session resources after any active command exits.
func (s *Session) Close() error {
	if s == nil {
		return nil
	}

	s.lifecycleMu.Lock()
	s.closed.Store(true)
	s.finishBreakpointsLocked("closed")
	s.lifecycleMu.Unlock()
	s.requestTermination()

	s.commandMu.Lock()
	defer s.commandMu.Unlock()
	s.resetValueReferences()
	s.hitBreakpointIDs = nil

	s.closeOnce.Do(func() {
		if s.started.Load() {
			s.closeErr = errors.Join(s.closeErr, s.runAfterHooks(context.Canceled))
		}

		if s.execution != nil && !s.executionClosed {
			s.executionClosed = true
			s.closeErr = errors.Join(s.closeErr, s.execution.Close())
		}

		if s.services != nil {
			s.closeErr = errors.Join(s.closeErr, s.services.Close())
		}
	})

	return s.closeErr
}

func (s *Session) checkContext(ctx context.Context) error {
	if ctx == nil {
		return runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	return ctx.Err()
}
