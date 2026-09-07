package ferret

import (
	"context"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// creationLifecycle orders parent cleanup after admitted constructors. It tracks
// in-flight creation only, never descendants. Callers own child cleanup.
// Transitions are protected by mu; callbacks and waiting happen without it.
// Closing done publishes the immutable cleanup result to concurrent closers.
type creationLifecycle struct {
	err      error
	idle     chan struct{}
	done     chan struct{}
	stopping chan struct{}
	active   int
	mu       sync.Mutex
	closing  bool
}

func (l *creationLifecycle) enter() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closing {
		return runtime.Error(runtime.ErrInvalidOperation, "parent is closed")
	}

	l.active++

	return nil
}

func (l *creationLifecycle) leave() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.active--
	if l.active == 0 && l.closing {
		close(l.idle)
	}
}

func (l *creationLifecycle) check() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closing {
		return runtime.Error(runtime.ErrInvalidOperation, "parent is closed")
	}

	return nil
}

func (l *creationLifecycle) close(cleanup func() error) error {
	l.mu.Lock()
	if l.closing {
		done := l.done
		l.mu.Unlock()
		<-done

		return l.err
	}

	l.closing = true
	if l.stopping != nil {
		close(l.stopping)
	}

	l.done = make(chan struct{})

	active := l.active
	if active != 0 {
		l.idle = make(chan struct{})
	}

	idle := l.idle
	l.mu.Unlock()

	if active != 0 {
		<-idle
	}

	err := cleanup()
	l.mu.Lock()
	l.err = err
	close(l.done)
	l.mu.Unlock()

	return err
}

// checkContext validates admission before executing user options/hooks.
func (l *creationLifecycle) checkContext(ctx context.Context) error {
	if ctx == nil {
		return runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return l.check()
}

// stopped supplies cancellation for a constructor blocked on parent capacity.
// Allocate the signal lazily so uncontended creation needs no channel allocation.
func (l *creationLifecycle) stopped() <-chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.stopping == nil {
		l.stopping = make(chan struct{})
		if l.closing {
			close(l.stopping)
		}
	}

	return l.stopping
}
