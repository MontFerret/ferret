package debugger

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/debugpoint"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Session controls one retained-state source-level debug execution.
// Command calls are serialized. Pause is safe to call concurrently with a
// running command.
type Session struct {
	execution        vm.DebugExecution
	values           vm.DebugValueAccess
	services         SessionServices
	closeErr         error
	runCtx           context.Context
	executionCtx     context.Context
	runCancel        context.CancelCauseFunc
	breakpoints      map[BreakpointID]Breakpoint
	boundPointIDs    map[BreakpointID]bytecode.DebugPointID
	source           *source.Source
	activeCancel     context.CancelCauseFunc
	valueRefs        map[ValueReference]runtime.Value
	params           []string
	pointIndex       debugpoint.Index
	format           FormatOptions
	nextBreakpointID BreakpointID
	nextValueRef     ValueReference
	activeCancelID   uint64
	closeOnce        sync.Once
	commandMu        sync.Mutex
	lifecycleMu      sync.Mutex
	closed           atomic.Bool
	started          atomic.Bool
	afterRun         bool
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

	if config.Source == nil {
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

	return &Session{
		execution:        config.Execution,
		values:           config.Values,
		services:         config.Services,
		source:           config.Source,
		pointIndex:       pointIndex,
		params:           append([]string(nil), config.Params...),
		format:           format,
		breakpoints:      make(map[BreakpointID]Breakpoint),
		boundPointIDs:    make(map[BreakpointID]bytecode.DebugPointID),
		nextBreakpointID: 1,
		valueRefs:        make(map[ValueReference]runtime.Value),
		nextValueRef:     1,
	}, nil
}

// Start begins execution and stops at the first executable source location.
func (s *Session) Start(ctx context.Context) (*Event, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}

	if s.started.Load() {
		return nil, &StateError{Operation: "start", State: "started"}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	runCtx, err := s.services.BeforeRun(ctx)
	if err != nil {
		return nil, fmt.Errorf("before run hooks: %w", err)
	}
	executionCtx, runCancel := context.WithCancelCause(runCtx)

	s.started.Store(true)
	s.runCtx = runCtx
	s.executionCtx = executionCtx
	s.setRunCancel(runCancel)
	s.resetValueReferences()

	event, err := s.execution.Start(s.services.ExtendContext(executionCtx))
	if err != nil {
		return nil, err
	}

	return s.convertEvent(event)
}

// Continue resumes execution until a breakpoint, pause request, error, or completion.
func (s *Session) Continue(ctx context.Context) (*Event, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	return s.resume(ctx, vm.DebugResumeContinue)
}

// Step stops at the next logical source location, including inside calls.
func (s *Session) Step(ctx context.Context) (*Event, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	return s.resume(ctx, vm.DebugResumeStep)
}

// Next stops at the next logical source location at the same or shallower call depth.
func (s *Session) Next(ctx context.Context) (*Event, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	return s.resume(ctx, vm.DebugResumeNext)
}

// Out stops at the next logical source location in a caller. At main it runs to completion.
func (s *Session) Out(ctx context.Context) (*Event, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	return s.resume(ctx, vm.DebugResumeOut)
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

// SetBreakpoint adds a source-line breakpoint using the legacy friendly
// next-executable-in-file binding policy.
func (s *Session) SetBreakpoint(file string, line int) (Breakpoint, error) {
	return s.SetBreakpointAt(
		SourceLocation{File: file, Line: line},
		BreakpointOptions{BindingMode: BreakpointBindNextExecutableInFile},
	)
}

// SetBreakpointAt adds a breakpoint at an explicit source location.
func (s *Session) SetBreakpointAt(location SourceLocation, opts BreakpointOptions) (Breakpoint, error) {
	if err := s.lockCommand(); err != nil {
		return Breakpoint{}, err
	}
	defer s.commandMu.Unlock()

	if err := s.ensureOpen(); err != nil {
		return Breakpoint{}, err
	}

	if location.Line <= 0 {
		return Breakpoint{}, runtime.Error(runtime.ErrInvalidArgument, "breakpoint line must be positive")
	}
	if location.Column < 0 {
		return Breakpoint{}, runtime.Error(runtime.ErrInvalidArgument, "breakpoint column must not be negative")
	}
	if opts.BindingMode < BreakpointBindNextExecutableInFile || opts.BindingMode > BreakpointBindNextExecutableInFunction {
		return Breakpoint{}, runtime.Errorf(runtime.ErrInvalidArgument, "unknown breakpoint binding mode %d", opts.BindingMode)
	}

	if location.File == "" {
		location.File = s.source.Name()
	}

	breakpoint := Breakpoint{
		ID:              s.nextBreakpointID,
		File:            location.File,
		RequestedLine:   location.Line,
		RequestedColumn: location.Column,
		BindingMode:     opts.BindingMode,
	}
	s.nextBreakpointID++

	if location.File == s.source.Name() {
		if point := s.breakpointPoint(location, opts.BindingMode); point != nil {
			bestLine, bestColumn := s.source.LocationAt(point.Span)
			breakpoint.Bound = true
			breakpoint.PointID = point.ID
			breakpoint.FunctionID = point.FunctionID
			breakpoint.Line = bestLine
			breakpoint.Column = bestColumn
			s.boundPointIDs[breakpoint.ID] = point.ID
		}
	}

	s.breakpoints[breakpoint.ID] = breakpoint

	return breakpoint, nil
}

// DeleteBreakpoint removes a breakpoint by ID.
func (s *Session) DeleteBreakpoint(id BreakpointID) error {
	if err := s.lockCommand(); err != nil {
		return err
	}
	defer s.commandMu.Unlock()

	if err := s.ensureOpen(); err != nil {
		return err
	}

	if _, exists := s.breakpoints[id]; !exists {
		return runtime.Errorf(runtime.ErrNotFound, "breakpoint %d", id)
	}

	delete(s.breakpoints, id)
	delete(s.boundPointIDs, id)

	return nil
}

// Breakpoints returns a stable ID-ordered snapshot.
func (s *Session) Breakpoints() []Breakpoint {
	if s == nil {
		return nil
	}

	s.commandMu.Lock()
	defer s.commandMu.Unlock()

	out := make([]Breakpoint, 0, len(s.breakpoints))
	for _, breakpoint := range s.breakpoints {
		out = append(out, breakpoint)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })

	return out
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
			FunctionID: frame.FunctionID,
			Location:   s.locationForPC(frame.PC, frame.FunctionID),
		})
	}

	return out, nil
}

// Locals returns the visible top-frame locals followed by bound parameters.
func (s *Session) Locals() ([]Variable, error) {
	if err := s.lockCommand(); err != nil {
		return nil, err
	}
	defer s.commandMu.Unlock()

	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	locals, err := s.execution.Locals()

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
	if err := s.lockCommand(); err != nil {
		return Value{}, err
	}
	defer s.commandMu.Unlock()

	if err := s.ensureOpen(); err != nil {
		return Value{}, err
	}
	locals, err := s.execution.Locals()

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

	if ctx == nil || ctx == s.runCtx || ctx == s.executionCtx {
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

	event, err := s.execution.Resume(s.services.ExtendContext(ctx), mode, s.breakpointPCs())
	if err != nil {
		return nil, err
	}

	return s.convertEvent(event)
}

func (s *Session) breakpointPCs() map[int]struct{} {
	out := make(map[int]struct{}, len(s.boundPointIDs))

	for _, pointID := range s.boundPointIDs {
		if point := s.pointIndex.PointByID(pointID); point != nil {
			out[point.PC] = struct{}{}
		}
	}

	return out
}

func (s *Session) breakpointPoint(location SourceLocation, mode BreakpointBindingMode) *bytecode.DebugPoint {
	exact := s.exactBreakpointPoint(location)
	if exact != nil || mode == BreakpointBindExact {
		return exact
	}

	if mode == BreakpointBindNextExecutableInFunction {
		return s.nextBreakpointPointInFunction(location)
	}

	return s.nextBreakpointPoint(location)
}

func (s *Session) exactBreakpointPoint(location SourceLocation) *bytecode.DebugPoint {
	var best *bytecode.DebugPoint
	points := s.pointIndex.Points()

	for i := range points {
		point := &points[i]
		pointLine, pointColumn := s.source.LocationAt(point.Span)

		if pointLine != location.Line || (location.Column > 0 && pointColumn != location.Column) {
			continue
		}

		if best == nil || s.debugPointLess(point, best) {
			best = point
		}
	}

	return best
}

func (s *Session) nextBreakpointPoint(location SourceLocation) *bytecode.DebugPoint {
	var best *bytecode.DebugPoint
	points := s.pointIndex.Points()

	for i := range points {
		point := &points[i]
		line, column := s.source.LocationAt(point.Span)
		if sourcePositionBefore(line, column, location.Line, location.Column) {
			continue
		}
		if best == nil || s.debugPointLess(point, best) {
			best = point
		}
	}

	return best
}

func (s *Session) nextBreakpointPointInFunction(location SourceLocation) *bytecode.DebugPoint {
	var previous, next *bytecode.DebugPoint
	previousAmbiguous := false
	nextAmbiguous := false
	points := s.pointIndex.Points()

	for i := range points {
		point := &points[i]
		line, column := s.source.LocationAt(point.Span)
		if sourcePositionBefore(line, column, location.Line, location.Column) {
			switch {
			case previous == nil || sourcePointPositionBefore(s.source, previous, point):
				previous = point
				previousAmbiguous = false
			case sameSourcePosition(s.source, previous, point) && previous.FunctionID != point.FunctionID:
				previousAmbiguous = true
			}
			continue
		}

		switch {
		case next == nil || sourcePointPositionBefore(s.source, point, next):
			next = point
			nextAmbiguous = false
		case sameSourcePosition(s.source, next, point) && next.FunctionID != point.FunctionID:
			nextAmbiguous = true
		}
	}

	if previous == nil || next == nil || previousAmbiguous || nextAmbiguous || previous.FunctionID != next.FunctionID {
		return nil
	}

	return next
}

func (s *Session) debugPointLess(left, right *bytecode.DebugPoint) bool {
	leftLine, leftColumn := s.source.LocationAt(left.Span)
	rightLine, rightColumn := s.source.LocationAt(right.Span)

	if leftLine != rightLine {
		return leftLine < rightLine
	}
	if leftColumn != rightColumn {
		return leftColumn < rightColumn
	}
	if left.PC != right.PC {
		return left.PC < right.PC
	}

	return left.ID < right.ID
}

func (s *Session) convertEvent(event *vm.DebugExecutionEvent) (*Event, error) {
	if event == nil {
		return nil, runtime.Error(runtime.ErrUnexpected, "debug execution returned no event")
	}
	out := &Event{Depth: event.Depth, Error: event.Error}

	if event.Point != nil {
		out.Location = s.location(event.Point.Span)
	}

	if s.closed.Load() && event.Reason != vm.DebugStopCompleted && event.Reason != vm.DebugStopTerminated {
		out.Reason = ReasonTerminated
		out.Error = context.Cause(s.executionCtx)
		if out.Error == nil {
			out.Error = context.Canceled
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
		if event.Point != nil {
			for id, pointID := range s.boundPointIDs {
				if pointID == event.Point.ID {
					out.HitBreakpointIDs = append(out.HitBreakpointIDs, id)
				}
			}
			sort.Slice(out.HitBreakpointIDs, func(i, j int) bool {
				return out.HitBreakpointIDs[i] < out.HitBreakpointIDs[j]
			})
		}
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
		if outputErr != nil || closeErr != nil || hookErr != nil {
			return nil, errors.Join(outputErr, closeErr, hookErr)
		}
		out.Output = output
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

func (s *Session) location(span source.Span) Location {
	line, column := s.source.LocationAt(span)

	return Location{File: s.source.Name(), Line: line, Column: column, Span: span}
}

func (s *Session) locationForPC(pc, functionID int) Location {
	var point *bytecode.DebugPoint

	if functionID >= -1 {
		point = s.pointIndex.NearestBeforeOrAtInFunction(functionID, pc)
	} else {
		point = s.pointIndex.NearestBeforeOrAt(pc)
	}

	if point == nil {
		return Location{File: s.source.Name()}
	}

	return s.location(point.Span)
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

	s.closed.Store(true)
	s.requestTermination()

	s.commandMu.Lock()
	defer s.commandMu.Unlock()
	s.resetValueReferences()

	s.closeOnce.Do(func() {
		if s.started.Load() {
			s.closeErr = errors.Join(s.closeErr, s.runAfterHooks(context.Canceled))
		}

		if s.execution != nil {
			s.closeErr = errors.Join(s.closeErr, s.execution.Close())
		}

		if s.services != nil {
			s.closeErr = errors.Join(s.closeErr, s.services.Close())
		}
	})

	return s.closeErr
}
