package debugger

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type fakeExecution struct {
	startEvent        *vm.DebugExecutionEvent
	resumeEvent       *vm.DebugExecutionEvent
	resumeBreakpoints map[int]struct{}
	locals            []vm.DebugLocal
	params            runtime.Params
	frames            []vm.DebugFrame
	status            vm.DebugExecutionStatus
	closed            bool
}

func (f *fakeExecution) Start(context.Context) (*vm.DebugExecutionEvent, error) {
	f.status = vm.DebugExecutionPaused
	return f.startEvent, nil
}

func (f *fakeExecution) Resume(_ context.Context, _ vm.DebugResumeMode, breakpoints map[int]struct{}) (*vm.DebugExecutionEvent, error) {
	f.resumeBreakpoints = breakpoints
	f.status = vm.DebugExecutionPaused
	return f.resumeEvent, nil
}

func (f *fakeExecution) RequestPause() {}

func (f *fakeExecution) Status() vm.DebugExecutionStatus {
	return f.status
}

func (f *fakeExecution) Locals() ([]vm.DebugLocal, error) {
	return f.locals, nil
}

func (f *fakeExecution) Params() runtime.Params {
	return f.params
}

func (f *fakeExecution) Frames() ([]vm.DebugFrame, error) {
	return f.frames, nil
}

func (f *fakeExecution) Close() error {
	f.closed = true
	return nil
}
