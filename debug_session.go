package ferret

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// DebugSession controls one retained-state source-level debug execution.
// It is not safe for concurrent command calls. Pause is safe to call
// concurrently with a running command.
type DebugSession struct {
	logger            logging.Logger
	closeErr          error
	runCtx            context.Context
	hooks             sessionHooks
	fs                fs.FileSystem
	limiter           *sessionLimiter
	encoding          *encoding.Registry
	execution         *vm.DebugExecution
	breakpoints       map[int]DebugBreakpoint
	boundBreakpointPC map[int]int
	source            *source.Source
	outputContentType string
	debugPoints       []bytecode.DebugPoint
	params            []string
	format            vm.DebugFormatOptions
	nextBreakpointID  int
	closeOnce         sync.Once
	closed            atomic.Bool
	started           atomic.Bool
	afterRun          bool
}

func newDebugSession(config debugSessionConfig) *DebugSession {
	return &DebugSession{
		execution:         config.execution,
		source:            config.source,
		debugPoints:       config.debugPoints,
		params:            append([]string(nil), config.params...),
		hooks:             config.hooks,
		limiter:           config.limiter,
		encoding:          config.encoding,
		logger:            config.logger,
		fs:                config.fs,
		outputContentType: config.outputContentType,
		format:            config.format,
		breakpoints:       make(map[int]DebugBreakpoint),
		boundBreakpointPC: make(map[int]int),
		nextBreakpointID:  1,
	}
}

// Start begins execution and stops at the first executable source location.
func (s *DebugSession) Start(ctx context.Context) (*DebugEvent, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	if s.started.Load() {
		return nil, &DebugStateError{Operation: "start", State: "started"}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	runCtx, err := s.hooks.runBeforeRunHooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("before run hooks: %w", err)
	}
	s.started.Store(true)
	s.runCtx = runCtx
	event, err := s.execution.Start(s.extendContext(runCtx))
	if err != nil {
		return nil, err
	}
	return s.convertEvent(event)
}

// Continue resumes execution until a breakpoint, pause request, error, or completion.
func (s *DebugSession) Continue(ctx context.Context) (*DebugEvent, error) {
	return s.resume(ctx, vm.DebugResumeContinue)
}

// Step stops at the next logical source location, including inside calls.
func (s *DebugSession) Step(ctx context.Context) (*DebugEvent, error) {
	return s.resume(ctx, vm.DebugResumeStep)
}

// Next stops at the next logical source location at the same or shallower call depth.
func (s *DebugSession) Next(ctx context.Context) (*DebugEvent, error) {
	return s.resume(ctx, vm.DebugResumeNext)
}

// Out stops at the next logical source location in a caller. At main it runs to completion.
func (s *DebugSession) Out(ctx context.Context) (*DebugEvent, error) {
	return s.resume(ctx, vm.DebugResumeOut)
}

// Pause requests a stop at the next logical source location.
func (s *DebugSession) Pause() error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if !s.started.Load() {
		return &DebugStateError{Operation: "pause", State: "new"}
	}
	s.execution.RequestPause()
	return nil
}

// SetBreakpoint adds a source-line breakpoint. When the requested line is not
// executable, it binds to the next executable line in the same source.
func (s *DebugSession) SetBreakpoint(file string, line int) (DebugBreakpoint, error) {
	if err := s.ensureOpen(); err != nil {
		return DebugBreakpoint{}, err
	}
	if line <= 0 {
		return DebugBreakpoint{}, runtime.Error(runtime.ErrInvalidArgument, "breakpoint line must be positive")
	}
	if file == "" && s.source != nil {
		file = s.source.Name()
	}

	breakpoint := DebugBreakpoint{ID: s.nextBreakpointID, File: file, RequestedLine: line}
	s.nextBreakpointID++
	if s.source != nil && file == s.source.Name() {
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
func (s *DebugSession) DeleteBreakpoint(id int) error {
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
func (s *DebugSession) Breakpoints() []DebugBreakpoint {
	if s == nil {
		return nil
	}
	out := make([]DebugBreakpoint, 0, len(s.breakpoints))
	for _, breakpoint := range s.breakpoints {
		out = append(out, breakpoint)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Frames returns the current frame followed by callers.
func (s *DebugSession) Frames() ([]DebugFrame, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	frames, err := s.execution.Frames()
	if err != nil {
		return nil, err
	}
	out := make([]DebugFrame, 0, len(frames))
	for _, frame := range frames {
		out = append(out, DebugFrame{
			Name:       frame.Name,
			FunctionID: frame.FunctionID,
			Location:   s.locationForPC(frame.PC, frame.FunctionID),
		})
	}
	return out, nil
}

// Locals returns the visible top-frame locals followed by bound parameters.
func (s *DebugSession) Locals() ([]DebugVariable, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	locals, err := s.execution.Locals()
	if err != nil {
		return nil, err
	}
	out := make([]DebugVariable, 0, len(locals)+len(s.execution.Params()))
	for _, local := range locals {
		out = append(out, DebugVariable{Name: local.Name, Mutable: local.Mutable, Value: s.debugValue(local.Value)})
	}
	params := s.execution.Params()
	names := append([]string(nil), s.params...)
	sort.Strings(names)
	for _, name := range names {
		if value, exists := params.Get(name); exists {
			out = append(out, DebugVariable{Name: "@" + name, Param: true, Value: s.debugValue(value)})
		}
	}
	return out, nil
}

// Evaluate evaluates a conservative, side-effect-free expression against the
// paused top frame.
func (s *DebugSession) Evaluate(ctx context.Context, expression string) (DebugValue, error) {
	if err := s.ensureOpen(); err != nil {
		return DebugValue{}, err
	}
	locals, err := s.execution.Locals()
	if err != nil {
		return DebugValue{}, err
	}
	values := make(map[string]runtime.Value, len(locals))
	for _, local := range locals {
		values[local.Name] = local.Value
	}
	value, err := evaluateDebugExpression(ctx, expression, values, s.execution.Params())
	if err != nil {
		return DebugValue{}, err
	}
	return s.debugValue(value), nil
}

func (s *DebugSession) resume(ctx context.Context, mode vm.DebugResumeMode) (*DebugEvent, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	if !s.started.Load() {
		return nil, &DebugStateError{Operation: "resume", State: "new"}
	}
	switch s.execution.Status() {
	case vm.DebugExecutionCompleted:
		return nil, &DebugStateError{Operation: "resume", State: "completed"}
	case vm.DebugExecutionTerminated:
		return nil, &DebugStateError{Operation: "resume", State: "terminated"}
	}
	if ctx == nil {
		ctx = s.runCtx
	}
	event, err := s.execution.Resume(s.extendContext(ctx), mode, s.breakpointPCs())
	if err != nil {
		return nil, err
	}
	return s.convertEvent(event)
}

func (s *DebugSession) breakpointPCs() map[int]struct{} {
	out := make(map[int]struct{}, len(s.boundBreakpointPC))
	for _, pc := range s.boundBreakpointPC {
		out[pc] = struct{}{}
	}
	return out
}

func (s *DebugSession) convertEvent(event *vm.DebugExecutionEvent) (*DebugEvent, error) {
	if event == nil {
		return nil, runtime.Error(runtime.ErrUnexpected, "debug execution returned no event")
	}
	out := &DebugEvent{Depth: event.Depth, Error: event.Error}
	if event.Point != nil {
		out.Location = s.location(event.Point.Span)
	}
	switch event.Reason {
	case vm.DebugStopEntry:
		out.Reason = DebugReasonEntry
	case vm.DebugStopBreakpoint:
		out.Reason = DebugReasonBreakpoint
	case vm.DebugStopStep:
		out.Reason = DebugReasonStep
	case vm.DebugStopPause:
		out.Reason = DebugReasonPause
	case vm.DebugStopRuntimeError:
		out.Reason = DebugReasonRuntimeError
		out.Error = errors.Join(out.Error, s.runAfterHooks(event.Error))
	case vm.DebugStopCompleted:
		out.Reason = DebugReasonCompleted
		output, outputErr := newOutput(s.encoding, s.outputContentType, event.Result)
		closeErr := event.Result.Close()
		hookErr := s.runAfterHooks(nil)
		if outputErr != nil || closeErr != nil || hookErr != nil {
			return nil, errors.Join(outputErr, closeErr, hookErr)
		}
		out.Output = output
	case vm.DebugStopTerminated:
		out.Reason = DebugReasonTerminated
		out.Error = errors.Join(out.Error, s.runAfterHooks(event.Error))
	}
	return out, nil
}

func (s *DebugSession) runAfterHooks(runErr error) error {
	if s.afterRun {
		return nil
	}
	s.afterRun = true
	if hookErr := s.hooks.runAfterRunHooks(s.runCtx, runErr); hookErr != nil {
		return fmt.Errorf("after run hooks: %w", hookErr)
	}
	return nil
}

func (s *DebugSession) debugValue(value runtime.Value) DebugValue {
	return DebugValue{
		Type:    vm.DebugValueTypeName(value),
		Display: vm.DebugFormatValue(value, s.format),
	}
}

func (s *DebugSession) location(span source.Span) DebugLocation {
	line, column := s.source.LocationAt(span)
	return DebugLocation{File: s.source.Name(), Line: line, Column: column, Span: span}
}

func (s *DebugSession) locationForPC(pc, functionID int) DebugLocation {
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
		return DebugLocation{File: s.source.Name()}
	}
	return s.location(found.Span)
}

func (s *DebugSession) extendContext(ctx context.Context) context.Context {
	ctx = s.logger.WithContext(ctx)
	ctx = encoding.WithRegistry(ctx, s.encoding)
	ctx = fs.WithFileSystem(ctx, s.fs)
	return ctx
}

func (s *DebugSession) ensureOpen() error {
	if s == nil || s.closed.Load() {
		return &DebugStateError{Operation: "use", State: "closed"}
	}
	return nil
}

// Close releases retained VM state, runs close hooks, and releases the plan's
// active-session permit.
func (s *DebugSession) Close() error {
	if s == nil {
		return nil
	}
	s.closeOnce.Do(func() {
		s.closed.Store(true)
		if s.execution != nil {
			s.closeErr = errors.Join(s.closeErr, s.execution.Close())
		}
		if s.hooks != nil {
			if err := s.hooks.runCloseHooks(); err != nil {
				s.closeErr = errors.Join(s.closeErr, fmt.Errorf("close hooks: %w", err))
			}
		}
		if s.limiter != nil {
			s.limiter.Release()
			s.limiter = nil
		}
	})
	return s.closeErr
}
