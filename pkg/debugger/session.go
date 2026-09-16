package debugger

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Session controls one retained-state source-level debug execution.
// Execution and inspection commands are serialized. Breakpoint mutation,
// breakpoint listing, Pause, and Close are safe during a running command.
// Context arguments must be non-nil. Inspection checks cancellation before and
// after command admission; cancellation does not interrupt a command-lock wait.
type Session struct {
	breakpoints breakpointSet
	inspector   inspector
	pointIndex  debugpoint.Index
	source      source.Source
	lifecycle   sessionLifecycle
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
		lifecycle:  newSessionLifecycle(config.Execution, config.Services),
		source:     config.Source,
		pointIndex: pointIndex,
	}
	session.inspector = newInspector(config, &session.source, &session.pointIndex, format)
	session.breakpoints.initialize(&session.lifecycle, &session.source, &session.pointIndex)

	return session, nil
}

// Start begins execution and stops at the first executable source location.
func (s *Session) Start(ctx context.Context) (*Event, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.lifecycle.unlockCommand()

	if err := s.lifecycle.prepareStart(ctx); err != nil {
		return nil, err
	}

	s.inspector.resetValueReferences()
	event, err := s.lifecycle.start()
	if err != nil {
		return nil, err
	}

	return s.convertEvent(event)
}

// Continue resumes execution until a breakpoint, pause request, error, or completion.
func (s *Session) Continue(ctx context.Context) (*Event, error) {
	return s.resume(ctx, vm.DebugResumeContinue)
}

// StepIn resumes execution until the next debuggable source location and may
// enter called functions.
func (s *Session) StepIn(ctx context.Context) (*Event, error) {
	return s.resume(ctx, vm.DebugResumeStepIn)
}

// StepOver resumes execution until the next debuggable source location at the
// same or shallower call depth. Breakpoints and pause requests may interrupt it
// inside deeper calls.
func (s *Session) StepOver(ctx context.Context) (*Event, error) {
	return s.resume(ctx, vm.DebugResumeStepOver)
}

// StepOut resumes execution until the current frame exits, then stops at the
// next debuggable source location in a caller. At main it runs to completion.
// Breakpoints and pause requests may interrupt it before the frame exits.
func (s *Session) StepOut(ctx context.Context) (*Event, error) {
	return s.resume(ctx, vm.DebugResumeStepOut)
}

// Pause requests a stop at the next logical source location.
func (s *Session) Pause(ctx context.Context) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	if err := s.ensureOpen(); err != nil {
		return err
	}

	return s.lifecycle.pause()
}

// ReplaceBreakpoints atomically replaces one source's complete breakpoint set.
// It is valid before execution, while paused, and while running. Empty requests
// clear the source; an empty sourceName selects the launched source. Results
// follow request order. Unbound results are successful verification outcomes;
// invalid positions or binding modes fail the whole operation.
//
// Unchanged requests retain IDs, including duplicate requests matched in ID
// order. Removed IDs are never reused. Other sources are unaffected. A stop
// already decided using a previous snapshot retains its hit IDs and remains
// inspectable; additions only affect subsequent source-point checks.
//
// Cancellation or native termination before publication leaves the set intact.
// Once publication commits, the operation returns success even if cancellation
// follows. Concurrent writers are ordered by admission to the writer gate.
func (s *Session) ReplaceBreakpoints(ctx context.Context, sourceName string, requests []BreakpointRequest) ([]Breakpoint, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}

	return s.breakpoints.replace(ctx, sourceName, requests)
}

// SetBreakpoint adds a location breakpoint using next-executable-in-source binding.
// It is safe while execution is running.
func (s *Session) SetBreakpoint(ctx context.Context, location source.Location) (Breakpoint, error) {
	return s.SetBreakpointAt(ctx, location, BreakpointOptions{BindingMode: BreakpointBindNextExecutableInSource})
}

// SetBreakpointAt adds a breakpoint at an explicit source location, including
// while execution is running. Each addition receives a distinct identity.
func (s *Session) SetBreakpointAt(ctx context.Context, location source.Location, opts BreakpointOptions) (Breakpoint, error) {
	if err := checkContext(ctx); err != nil {
		return Breakpoint{}, err
	}

	if err := s.ensureOpen(); err != nil {
		return Breakpoint{}, err
	}

	return s.breakpoints.add(ctx, location, opts)
}

// DeleteBreakpoint removes a breakpoint by ID, including while execution is
// running. A previously decided hit remains a valid stop.
func (s *Session) DeleteBreakpoint(ctx context.Context, id BreakpointID) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	if err := s.ensureOpen(); err != nil {
		return err
	}

	return s.breakpoints.delete(ctx, id)
}

// Breakpoints returns a detached, ID-ordered snapshot, including while running
// or after Close. Listing never waits for execution or an in-progress replacement.
func (s *Session) Breakpoints(ctx context.Context) ([]Breakpoint, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if s == nil {
		return nil, nil
	}

	return s.breakpoints.list(), nil
}

// Frames returns the current frame followed by callers.
func (s *Session) Frames(ctx context.Context) ([]Frame, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.lifecycle.unlockCommand()

	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}

	return s.inspector.frames()
}

// Locals returns the visible top-frame locals followed by bound parameters.
func (s *Session) Locals(ctx context.Context) ([]Variable, error) {
	return s.FrameLocals(ctx, 0)
}

// FrameLocals returns the visible locals and bound parameters for one paused
// frame. Frame indexes follow Frames: zero is the current frame.
func (s *Session) FrameLocals(ctx context.Context, frame int) ([]Variable, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.lifecycle.unlockCommand()

	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}

	return s.inspector.frameLocals(frame)
}

// Variables returns the child variables for one expandable debugger value from
// the current paused state.
func (s *Session) Variables(ctx context.Context, reference ValueReference) ([]Variable, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.lifecycle.unlockCommand()

	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}

	return s.inspector.variables(reference)
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
	if err := checkContext(ctx); err != nil {
		return Value{}, err
	}

	if err := s.lockCommand(); err != nil {
		return Value{}, err
	}
	defer s.lifecycle.unlockCommand()

	if err := checkContext(ctx); err != nil {
		return Value{}, err
	}

	if err := s.ensureOpen(); err != nil {
		return Value{}, err
	}

	return s.inspector.evaluateFrame(ctx, frame, expression)
}

// Close requests termination, then releases retained VM state and
// embedding-owned session resources after any active command exits.
func (s *Session) Close() error {
	if s == nil {
		return nil
	}

	s.lifecycle.beginClose()
	defer s.lifecycle.unlockCommand()
	s.inspector.resetValueReferences()
	s.breakpoints.resetHit()

	return s.lifecycle.finishClose()
}

func (s *Session) resume(ctx context.Context, mode vm.DebugResumeMode) (*Event, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.lifecycle.unlockCommand()

	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	command, err := s.lifecycle.prepareResume(ctx)
	if err != nil {
		return nil, err
	}

	if command.cleanup != nil {
		defer s.lifecycle.finishResume(command)
	}

	s.inspector.resetValueReferences()
	s.breakpoints.resetHit()
	event, err := s.lifecycle.resume(command.ctx, mode, s.breakpoints.predicate)
	if err != nil {
		return nil, err
	}

	return s.convertEvent(event)
}

func (s *Session) convertEvent(event *vm.DebugExecutionEvent) (*Event, error) {
	if event == nil {
		return nil, runtime.Error(runtime.ErrUnexpected, "debug execution returned no event")
	}

	if event.Reason == vm.DebugStopCompleted {
		s.lifecycle.finishTerminal("completed")
	} else if event.Reason == vm.DebugStopTerminated {
		s.lifecycle.finishTerminal("terminated")
	}

	out := &Event{Depth: event.Depth, Error: event.Error}
	if event.Point != nil {
		out.Location = s.source.RangeAt(event.Point.Span)
	}

	if s.lifecycle.closed.Load() && event.Reason != vm.DebugStopCompleted && event.Reason != vm.DebugStopTerminated {
		out.Reason = ReasonTerminated
		out.Error = s.lifecycle.terminateStop(event.Error)

		return out, nil
	}

	switch event.Reason {
	case vm.DebugStopEntry:
		out.Reason = ReasonEntry
	case vm.DebugStopBreakpoint:
		out.Reason = ReasonBreakpoint
		out.HitBreakpointIDs = s.breakpoints.hitBreakpointIDs()
	case vm.DebugStopStep:
		out.Reason = ReasonStep
	case vm.DebugStopPause:
		out.Reason = ReasonPause
	case vm.DebugStopRuntimeError:
		out.Reason = ReasonRuntimeError
		out.Error = s.lifecycle.settleRunError(event.Error)
	case vm.DebugStopCompleted:
		s.inspector.resetValueReferences()
		out.Reason = ReasonCompleted
		output, err := s.lifecycle.complete(event.Result)
		out.Output = output

		return out, err
	case vm.DebugStopTerminated:
		s.inspector.resetValueReferences()
		out.Reason = ReasonTerminated
		out.Error = s.lifecycle.settleRunError(event.Error)
	}

	return out, nil
}

func (s *Session) ensureOpen() error {
	if s == nil {
		return &StateError{Operation: "use", State: "closed"}
	}

	return s.lifecycle.ensureOpen()
}

func (s *Session) lockCommand() error {
	if s == nil {
		return &StateError{Operation: "use", State: "closed"}
	}

	return s.lifecycle.lockCommand()
}
