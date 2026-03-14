package ferret

import "context"

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

func (l *sessionLimiter) Acquire(ctx context.Context) error {
	if l == nil || l.ch == nil {
		return nil
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	select {
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
