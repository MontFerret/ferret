package debugger

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Session controls one retained-state source-level debug execution.
// It is not safe for concurrent command calls. Pause is safe to call
// concurrently with a running command.
type Session struct {
	closeErr          error
	runCtx            context.Context
	execution         vm.DebugExecution
	values            vm.DebugValueAccess
	services          SessionServices
	breakpoints       map[int]Breakpoint
	boundBreakpointPC map[int]int
	source            *source.Source
	debugPoints       []bytecode.DebugPoint
	params            []string
	format            FormatOptions
	nextBreakpointID  int
	closeOnce         sync.Once
	closed            atomic.Bool
	started           atomic.Bool
	afterRun          bool
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
	return &Session{
		execution:         config.Execution,
		values:            config.Values,
		services:          config.Services,
		source:            config.Source,
		debugPoints:       append([]bytecode.DebugPoint(nil), config.DebugPoints...),
		params:            append([]string(nil), config.Params...),
		format:            format,
		breakpoints:       make(map[int]Breakpoint),
		boundBreakpointPC: make(map[int]int),
		nextBreakpointID:  1,
	}, nil
}

// Start begins execution and stops at the first executable source location.
func (s *Session) Start(ctx context.Context) (*Event, error) {
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
	s.started.Store(true)
	s.runCtx = runCtx
	event, err := s.execution.Start(s.services.ExtendContext(runCtx))
	if err != nil {
		return nil, err
	}
	return s.convertEvent(event)
}

// Continue resumes execution until a breakpoint, pause request, error, or completion.
func (s *Session) Continue(ctx context.Context) (*Event, error) {
	return s.resume(ctx, vm.DebugResumeContinue)
}

// Step stops at the next logical source location, including inside calls.
func (s *Session) Step(ctx context.Context) (*Event, error) {
	return s.resume(ctx, vm.DebugResumeStep)
}

// Next stops at the next logical source location at the same or shallower call depth.
func (s *Session) Next(ctx context.Context) (*Event, error) {
	return s.resume(ctx, vm.DebugResumeNext)
}

// Out stops at the next logical source location in a caller. At main it runs to completion.
func (s *Session) Out(ctx context.Context) (*Event, error) {
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

// SetBreakpoint adds a source-line breakpoint. When the requested line is not
// executable, it binds to the next executable line in the same source.
func (s *Session) SetBreakpoint(file string, line int) (Breakpoint, error) {
	if err := s.ensureOpen(); err != nil {
		return Breakpoint{}, err
	}
	if line <= 0 {
		return Breakpoint{}, runtime.Error(runtime.ErrInvalidArgument, "breakpoint line must be positive")
	}
	if file == "" {
		file = s.source.Name()
	}

	breakpoint := Breakpoint{ID: s.nextBreakpointID, File: file, RequestedLine: line}
	s.nextBreakpointID++
	if file == s.source.Name() {
		bestLine := 0
		bestColumn := 0
		bestPC := -1
		for _, point := range s.debugPoints {
			pointLine, pointColumn := s.source.LocationAt(point.Span)
			if pointLine >= line &&
				(bestLine == 0 || pointLine < bestLine || (pointLine == bestLine && pointColumn < bestColumn)) {
				bestLine = pointLine
				bestColumn = pointColumn
				bestPC = point.PC
			}
		}
		if bestLine != 0 {
			breakpoint.Bound = true
			breakpoint.Line = bestLine
			breakpoint.Column = bestColumn
			s.boundBreakpointPC[breakpoint.ID] = bestPC
		}
	}
	s.breakpoints[breakpoint.ID] = breakpoint
	return breakpoint, nil
}

// DeleteBreakpoint removes a breakpoint by ID.
func (s *Session) DeleteBreakpoint(id int) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if _, exists := s.breakpoints[id]; !exists {
		return runtime.Errorf(runtime.ErrNotFound, "breakpoint %d", id)
	}
	delete(s.breakpoints, id)
	delete(s.boundBreakpointPC, id)
	return nil
}

// Breakpoints returns a stable ID-ordered snapshot.
func (s *Session) Breakpoints() []Breakpoint {
	if s == nil {
		return nil
	}
	out := make([]Breakpoint, 0, len(s.breakpoints))
	for _, breakpoint := range s.breakpoints {
		out = append(out, breakpoint)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Frames returns the current frame followed by callers.
func (s *Session) Frames() ([]Frame, error) {
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

// Evaluate evaluates a conservative, side-effect-free expression against the
// paused top frame.
func (s *Session) Evaluate(ctx context.Context, expression string) (Value, error) {
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

	if ctx == nil || ctx == s.runCtx {
		ctx = s.runCtx
	} else {
		var cleanup func()
		ctx, cleanup = newResumeContext(s.runCtx, ctx)
		defer cleanup()
	}

	event, err := s.execution.Resume(s.services.ExtendContext(ctx), mode, s.breakpointPCs())
	if err != nil {
		return nil, err
	}

	return s.convertEvent(event)
}

func (s *Session) breakpointPCs() map[int]struct{} {
	out := make(map[int]struct{}, len(s.boundBreakpointPC))
	for _, pc := range s.boundBreakpointPC {
		out[pc] = struct{}{}
	}
	return out
}

func (s *Session) convertEvent(event *vm.DebugExecutionEvent) (*Event, error) {
	if event == nil {
		return nil, runtime.Error(runtime.ErrUnexpected, "debug execution returned no event")
	}
	out := &Event{Depth: event.Depth, Error: event.Error}
	if event.Point != nil {
		out.Location = s.location(event.Point.Span)
	}
	switch event.Reason {
	case vm.DebugStopEntry:
		out.Reason = ReasonEntry
	case vm.DebugStopBreakpoint:
		out.Reason = ReasonBreakpoint
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
		out.Reason = ReasonCompleted
		output, outputErr := s.services.Materialize(event.Result)
		closeErr := event.Result.Close()
		hookErr := s.runAfterHooks(nil)
		if outputErr != nil || closeErr != nil || hookErr != nil {
			return nil, errors.Join(outputErr, closeErr, hookErr)
		}
		out.Output = output
	case vm.DebugStopTerminated:
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
	if hookErr := s.services.AfterRun(s.runCtx, runErr); hookErr != nil {
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
		Type:    typeName,
		Display: formatValueWithInfo(value, info, s.values, s.format),
	}
}

func (s *Session) location(span source.Span) Location {
	line, column := s.source.LocationAt(span)
	return Location{File: s.source.Name(), Line: line, Column: column, Span: span}
}

func (s *Session) locationForPC(pc, functionID int) Location {
	var found *bytecode.DebugPoint
	for i := range s.debugPoints {
		point := &s.debugPoints[i]
		if point.PC > pc {
			break
		}
		if point.FunctionID == functionID {
			found = point
		}
	}
	if found == nil {
		return Location{File: s.source.Name()}
	}
	return s.location(found.Span)
}

func (s *Session) ensureOpen() error {
	if s == nil || s.closed.Load() {
		return &StateError{Operation: "use", State: "closed"}
	}
	return nil
}

// Close releases retained VM state and embedding-owned session resources.
func (s *Session) Close() error {
	if s == nil {
		return nil
	}
	s.closeOnce.Do(func() {
		s.closed.Store(true)
		if s.execution != nil {
			s.closeErr = errors.Join(s.closeErr, s.execution.Close())
		}
		if s.services != nil {
			s.closeErr = errors.Join(s.closeErr, s.services.Close())
		}
	})
	return s.closeErr
}
