package vm

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// DebugStopReason identifies why an incremental VM execution stopped.
	DebugStopReason uint8
	// DebugResumeMode selects the stepping behavior for a resumed execution.
	DebugResumeMode uint8
	// DebugExecutionStatus reports the incremental execution lifecycle state.
	DebugExecutionStatus uint8

	// DebugExecutionEvent reports one incremental VM stop or terminal event.
	DebugExecutionEvent struct {
		Point  *bytecode.DebugPoint
		Result *Result
		Error  error
		Reason DebugStopReason
		Depth  int
	}

	// DebugLocal exposes one VM binding visible at the current source location.
	DebugLocal struct {
		Value   runtime.Value
		Name    string
		Mutable bool
	}

	// DebugFrame describes one retained VM frame.
	DebugFrame struct {
		Name       string
		FunctionID int
		PC         int
	}

	// DebugExecution controls one retained-state execution of a debug-enabled
	// program. Resume calls must be serialized. RequestPause and Close are safe to
	// call concurrently with execution.
	DebugExecution interface {
		Start(context.Context) (*DebugExecutionEvent, error)
		Resume(context.Context, DebugResumeMode, map[int]struct{}) (*DebugExecutionEvent, error)
		RequestPause()
		Status() DebugExecutionStatus
		Locals() ([]DebugLocal, error)
		Params() runtime.Params
		Frames() ([]DebugFrame, error)
		Close() error
	}

	debugExecution struct {
		terminalErr    error
		vm             *VM
		env            *Environment
		current        *bytecode.DebugPoint
		control        debugControl
		mu             sync.Mutex
		pauseRequested atomic.Bool
		status         DebugExecutionStatus
	}
)

const debugTrapOpcode = bytecode.Opcode(255)

const (
	DebugStopEntry DebugStopReason = iota
	DebugStopBreakpoint
	DebugStopStep
	DebugStopPause
	DebugStopRuntimeError
	DebugStopCompleted
	DebugStopTerminated
)

const (
	DebugResumeContinue DebugResumeMode = iota
	DebugResumeStep
	DebugResumeNext
	DebugResumeOut
)

const (
	DebugExecutionNew DebugExecutionStatus = iota
	DebugExecutionPaused
	DebugExecutionRunning
	DebugExecutionCompleted
	DebugExecutionTerminated
	DebugExecutionClosed
)

// NewDebugExecution creates an incremental execution for a program compiled
// with debugger metadata.
func NewDebugExecution(instance *VM, env *Environment) (DebugExecution, error) {
	if instance == nil || instance.closed {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "vm is closed")
	}

	if len(instance.program.Metadata.DebugPoints) == 0 {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "program has no debugger metadata")
	}

	if err := bytecode.ValidateProgram(instance.program); err != nil {
		return nil, err
	}

	if env == nil {
		env = noopEnv
	}

	exec := &debugExecution{vm: instance, env: env, status: DebugExecutionNew}
	exec.control = debugControl{
		owner:  exec,
		points: make(map[int]*bytecode.DebugPoint, len(instance.program.Metadata.DebugPoints)),
	}

	for i := range instance.program.Metadata.DebugPoints {
		point := &instance.program.Metadata.DebugPoints[i]

		if instance.plan.instructions[point.PC].Opcode != bytecode.OpJump {
			return nil, runtime.Error(runtime.ErrInvalidOperation, "vm is already configured for debugging")
		}

		exec.control.points[point.PC] = point
	}

	for i := range instance.program.Metadata.DebugPoints {
		point := &instance.program.Metadata.DebugPoints[i]
		instance.plan.instructions[point.PC].Opcode = debugTrapOpcode
	}

	return exec, nil
}

// Start begins execution and stops at the first executable source location.
func (d *debugExecution) Start(ctx context.Context) (*DebugExecutionEvent, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.status != DebugExecutionNew {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "debug execution has already started")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if err := d.vm.state.startRun(d.env); err != nil {
		d.status = DebugExecutionTerminated
		return nil, err
	}

	if err := warmup(d.vm, d.env); err != nil {
		d.terminalErr = d.vm.state.wrapRuntimeError(err)
		d.status = DebugExecutionPaused

		return &DebugExecutionEvent{Reason: DebugStopRuntimeError, Error: d.terminalErr}, nil
	}

	d.control.entry = true
	return d.runLocked(ctx)
}

// Resume continues a paused execution according to mode and the active
// breakpoint PCs.
func (d *debugExecution) Resume(ctx context.Context, mode DebugResumeMode, breakpoints map[int]struct{}) (*DebugExecutionEvent, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.status != DebugExecutionPaused {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "debug execution is not paused")
	}

	if d.terminalErr != nil {
		err := d.terminalErr
		d.vm.state.endRun()
		d.status = DebugExecutionTerminated

		return &DebugExecutionEvent{Reason: DebugStopTerminated, Error: err}, nil
	}

	if ctx == nil {
		ctx = context.Background()
	}

	d.control.mode = mode
	d.control.breakpoints = breakpoints
	d.control.startDepth = d.vm.state.frames.Len()

	if d.current != nil {
		d.control.skip = true
		d.control.skipPC = d.current.PC
		d.control.skipDepth = d.control.startDepth
	}

	return d.runLocked(ctx)
}

// RequestPause asks a running execution to pause at its next debug point.
func (d *debugExecution) RequestPause() {
	if d != nil {
		d.pauseRequested.Store(true)
	}
}

func (d *debugExecution) Status() DebugExecutionStatus {
	if d == nil {
		return DebugExecutionClosed
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	return d.status
}

// Locals returns values for bindings visible at the current stop.
func (d *debugExecution) Locals() ([]DebugLocal, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.status != DebugExecutionPaused || d.current == nil {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "debug execution is not paused at a source location")
	}

	out := make([]DebugLocal, 0, len(d.current.Bindings))
	for _, binding := range d.current.Bindings {
		value := d.vm.state.valueOf(d.vm.program.Constants, binding.Register)

		if binding.Cell {
			if handle, ok := d.vm.state.cellHandleOf(binding.Register); ok {
				if cellValue, exists := d.vm.state.cells.Get(handle); exists {
					value = cellValue
				}
			}
		}

		out = append(out, DebugLocal{Name: binding.Name, Value: value, Mutable: binding.Mutable})
	}

	return out, nil
}

// Params returns the bound host parameters for the current execution.
func (d *debugExecution) Params() runtime.Params {
	d.mu.Lock()
	defer d.mu.Unlock()

	out := runtime.NewParams()

	if d.env != nil {
		for name, value := range d.env.Params {
			out[name] = value
		}
	}

	return out
}

// Frames returns the current frame followed by callers from nearest to farthest.
func (d *debugExecution) Frames() ([]DebugFrame, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.status != DebugExecutionPaused {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "debug execution is not paused")
	}

	currentID := -1
	currentPC := d.vm.state.pc
	name := "<main>"

	if d.current != nil {
		currentID = d.current.FunctionID
		currentPC = d.current.PC

		if currentID >= 0 && currentID < len(d.vm.program.Functions.UserDefined) {
			name = d.vm.program.Functions.UserDefined[currentID].DisplayName

			if name == "" {
				name = d.vm.program.Functions.UserDefined[currentID].Name
			}
		}
	}

	out := []DebugFrame{{Name: name, FunctionID: currentID, PC: currentPC}}
	for _, trace := range d.vm.state.frames.TraceEntries() {
		callerID := -1
		callerName := "<main>"

		if point := d.nearestPoint(trace.CallSitePC); point != nil {
			callerID = point.FunctionID

			if callerID >= 0 && callerID < len(d.vm.program.Functions.UserDefined) {
				callerName = d.vm.program.Functions.UserDefined[callerID].DisplayName

				if callerName == "" {
					callerName = d.vm.program.Functions.UserDefined[callerID].Name
				}
			}
		}

		out = append(out, DebugFrame{Name: callerName, FunctionID: callerID, PC: trace.CallSitePC})
	}

	return out, nil
}

func (d *debugExecution) runLocked(ctx context.Context) (event *DebugExecutionEvent, err error) {
	d.status = DebugExecutionRunning

	defer func() {
		if recovered := recover(); recovered != nil {
			d.terminalErr = d.vm.state.runtimeErrorFromPanic(recovered)
			d.status = DebugExecutionPaused
			event = &DebugExecutionEvent{Reason: DebugStopRuntimeError, Error: d.terminalErr, Depth: d.vm.state.frames.Len(), Point: d.current}
			err = nil
		}
	}()

	for {
		root, runErr := d.vm.runCore(ctx, nil, true)

		if runErr != nil {
			if pc, ok := d.debugTrapPC(); ok {
				if d.control.shouldStop(pc, d.vm.state.frames.Len()) {
					d.vm.state.pc = pc
					d.status = DebugExecutionPaused
					d.current = d.control.points[pc]
					return &DebugExecutionEvent{Reason: d.control.reason, Point: d.current, Depth: d.vm.state.frames.Len()}, nil
				}

				continue
			}

			d.terminalErr = d.vm.state.wrapRuntimeError(runErr)
			d.status = DebugExecutionPaused
			d.current = d.nearestPoint(d.vm.state.lastPC)

			return &DebugExecutionEvent{Reason: DebugStopRuntimeError, Error: d.terminalErr, Point: d.current, Depth: d.vm.state.frames.Len()}, nil
		}

		result := d.vm.state.finishRun(root)
		d.status = DebugExecutionCompleted
		d.current = nil

		return &DebugExecutionEvent{Reason: DebugStopCompleted, Result: result}, nil
	}
}

func (d *debugExecution) debugTrapPC() (int, bool) {
	pc := d.vm.state.pc - 1

	if pc < 0 || pc >= len(d.vm.plan.instructions) {
		return 0, false
	}

	if d.vm.plan.instructions[pc].Opcode != debugTrapOpcode || d.control.points[pc] == nil {
		return 0, false
	}

	return pc, true
}

func (d *debugExecution) nearestPoint(pc int) *bytecode.DebugPoint {
	var found *bytecode.DebugPoint

	for i := range d.vm.program.Metadata.DebugPoints {
		point := &d.vm.program.Metadata.DebugPoints[i]

		if point.PC > pc {
			break
		}

		found = point
	}

	return found
}

// Close releases retained execution state and permanently closes the VM.
func (d *debugExecution) Close() error {
	if d == nil {
		return nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.status == DebugExecutionClosed {
		return nil
	}

	if d.vm != nil {
		d.vm.state.endRun()
		_ = d.vm.Close()
	}

	d.status = DebugExecutionClosed

	return nil
}
