package utils

import (
	"context"
	"sync"
)

// Signal when Wait reaches its cancellation select, without scheduler sleeps.
type observedWaitContext struct {
	context.Context
	ready chan struct{}
	once  sync.Once
}

func (ctx *observedWaitContext) Done() <-chan struct{} {
	ctx.once.Do(func() { close(ctx.ready) })

	return ctx.Context.Done()
}
