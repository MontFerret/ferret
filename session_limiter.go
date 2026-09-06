package ferret

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

func (l *sessionLimiter) Acquire(ctx context.Context, parents ...*creationLifecycle) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if l == nil || l.ch == nil {
		return nil
	}

	select {
	case l.ch <- struct{}{}:
		return nil
	default:
	}

	var planStopped, engineStopped <-chan struct{}

	if len(parents) > 0 && parents[0] != nil {
		planStopped = parents[0].stopped()
	}

	if len(parents) > 1 && parents[1] != nil {
		engineStopped = parents[1].stopped()
	}

	select {
	case <-planStopped:
		return runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	case <-engineStopped:
		return runtime.Error(runtime.ErrInvalidOperation, "engine is closed")
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
