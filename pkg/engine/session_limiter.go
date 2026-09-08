package engine

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type sessionLimiter struct {
	ch chan struct{}
}

func newSessionLimiter(max int) *sessionLimiter {
	if max <= 0 {
		return &sessionLimiter{}
	}

	return &sessionLimiter{
		ch: make(chan struct{}, max),
	}
}

func (l *sessionLimiter) Acquire(ctx context.Context, planClosed <-chan struct{}) error {
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

func (l *sessionLimiter) Release() {
	if l == nil || l.ch == nil {
		return
	}

	select {
	case <-l.ch:
	default:
	}
}
