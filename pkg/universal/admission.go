package universal

import (
	"context"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// admission tracks calls, never descendants. The mutex protects Add versus
// closure, so no Add can race a Wait that observes an empty group.
type admission struct {
	mu     sync.Mutex
	active sync.WaitGroup
	closed bool
}

func (a *admission) begin(ctx context.Context, owner string) error {
	if ctx == nil {
		return runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return runtime.Error(runtime.ErrInvalidOperation, owner+" is closed")
	}

	a.active.Add(1)

	return nil
}

func (a *admission) end() {
	a.active.Done()
}

func (a *admission) close() {
	a.mu.Lock()
	a.closed = true
	a.mu.Unlock()
}

func (a *admission) wait() {
	a.active.Wait()
}
