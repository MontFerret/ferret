package debugger

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type fakeExecution struct {
	closeErr          error
	startErr          error
	resumeErr         error
	startEvent        *vm.DebugExecutionEvent
	resumeEvent       *vm.DebugExecutionEvent
	resumeBreakpoints vm.DebugBreakpointPredicate
	params            runtime.Params
	frames            []vm.DebugFrame
	locals            []vm.DebugLocal
	closeCalls        int
	localsCalls       int
	paramsCalls       int
	framesCalls       int
	status            vm.DebugExecutionStatus
	errorStatus       vm.DebugExecutionStatus
	closed            bool
}

func (f *fakeExecution) Start(context.Context) (*vm.DebugExecutionEvent, error) {
	if f.startErr != nil {
		f.status = f.errorStatus

		return nil, f.startErr
	}

	f.status = vm.DebugExecutionPaused
	return f.startEvent, nil
}

func (f *fakeExecution) Resume(_ context.Context, _ vm.DebugResumeMode, breakpoints vm.DebugBreakpointPredicate) (*vm.DebugExecutionEvent, error) {
	if f.resumeErr != nil {
		f.status = f.errorStatus

		return nil, f.resumeErr
	}

	f.resumeBreakpoints = breakpoints
	f.status = vm.DebugExecutionPaused
	if f.resumeEvent != nil && f.resumeEvent.Reason == vm.DebugStopBreakpoint && f.resumeEvent.Point != nil && breakpoints != nil {
		breakpoints(f.resumeEvent.Point.PC)
	}
	return f.resumeEvent, nil
}

func (f *fakeExecution) RequestPause() {}

func (f *fakeExecution) Status() vm.DebugExecutionStatus {
	return f.status
}

func (f *fakeExecution) Locals() ([]vm.DebugLocal, error) {
	f.localsCalls++

	return f.locals, nil
}

func (f *fakeExecution) Params() runtime.Params {
	f.paramsCalls++

	return f.params
}

func (f *fakeExecution) Frames() ([]vm.DebugFrame, error) {
	f.framesCalls++

	return f.frames, nil
}

func (f *fakeExecution) Close() error {
	f.closeCalls++
	f.closed = true
	return f.closeErr
}
