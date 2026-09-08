// Package session supplies native session admission, debugger services, and output materialization.
package session

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Limiter bounds active sessions across plans without owning their lifetime.
type Limiter struct {
	ch chan struct{}
}

// NewLimiter disables the capacity bound when max is non-positive.
func NewLimiter(max int) *Limiter {
	if max <= 0 {
		return &Limiter{}
	}

	return &Limiter{
		ch: make(chan struct{}, max),
	}
}

// Acquire waits for capacity, observing the request context and Plan closure only.
// The caller must supply a non-nil context and release every acquired permit.
func (l *Limiter) Acquire(ctx context.Context, planClosed <-chan struct{}) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	select {
	case <-planClosed:
		return runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	default:
	}

	if l == nil || l.ch == nil {
		return nil
	}

	select {
	case <-planClosed:
		return runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	case l.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release returns one permit, if any. Session owners provide once-only release.
func (l *Limiter) Release() {
	if l == nil || l.ch == nil {
		return
	}

	select {
	case <-l.ch:
	default:
	}
}
