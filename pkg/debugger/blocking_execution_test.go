package debugger

import (
	"context"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type blockingExecution struct {
	resumeStarted  chan struct{}
	resumeRelease  chan struct{}
	resumeCanceled chan struct{}
	pauseCalled    chan struct{}
	point          bytecode.DebugPoint
	resumeCalls    int
	active         int
	maxActive      int
	closeCalls     int
	releaseOnce    sync.Once
	mu             sync.Mutex
	status         vm.DebugExecutionStatus
}

func newBlockingExecution(point bytecode.DebugPoint) *blockingExecution {
	return &blockingExecution{
		point:          point,
		resumeStarted:  make(chan struct{}, 4),
		resumeRelease:  make(chan struct{}),
		resumeCanceled: make(chan struct{}, 1),
		pauseCalled:    make(chan struct{}, 1),
		status:         vm.DebugExecutionNew,
	}
}

func (b *blockingExecution) Start(context.Context) (*vm.DebugExecutionEvent, error) {
	b.mu.Lock()
	b.status = vm.DebugExecutionPaused
	b.mu.Unlock()

	return &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &b.point}, nil
}

func (b *blockingExecution) Resume(ctx context.Context, _ vm.DebugResumeMode, _ map[int]struct{}) (*vm.DebugExecutionEvent, error) {
	b.mu.Lock()
	b.status = vm.DebugExecutionRunning
	b.resumeCalls++
	b.active++
	if b.active > b.maxActive {
		b.maxActive = b.active
	}
	b.mu.Unlock()

	b.resumeStarted <- struct{}{}
	<-b.resumeRelease
	if ctx.Err() != nil {
		select {
		case b.resumeCanceled <- struct{}{}:
		default:
		}
	}

	b.mu.Lock()
	b.active--
	b.status = vm.DebugExecutionPaused
	b.mu.Unlock()

	return &vm.DebugExecutionEvent{Reason: vm.DebugStopStep, Point: &b.point}, nil
}

func (b *blockingExecution) RequestPause() {
	select {
	case b.pauseCalled <- struct{}{}:
	default:
	}

	b.release()
}

func (b *blockingExecution) Status() vm.DebugExecutionStatus {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.status
}

func (b *blockingExecution) Locals() ([]vm.DebugLocal, error) {
	return nil, nil
}

func (b *blockingExecution) Params() runtime.Params {
	return runtime.NewParams()
}

func (b *blockingExecution) Frames() ([]vm.DebugFrame, error) {
	return nil, nil
}

func (b *blockingExecution) Close() error {
	b.mu.Lock()
	b.status = vm.DebugExecutionClosed
	b.closeCalls++
	b.mu.Unlock()
	b.release()

	return nil
}

func (b *blockingExecution) release() {
	b.releaseOnce.Do(func() {
		close(b.resumeRelease)
	})
}

func (b *blockingExecution) stats() (resumeCalls, maxActive, closeCalls int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.resumeCalls, b.maxActive, b.closeCalls
}
