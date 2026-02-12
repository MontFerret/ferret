package runtime

import (
	"context"
	"time"
)

type (
	// Context is an extension of the standard context.Context interface that includes additional methods for logging and memory allocation.
	Context interface {
		context.Context

		Logger() Logger
		Alloc() Allocator
	}

	contextImpl struct {
		base      context.Context
		logger    Logger
		allocator Allocator
	}
)

func NewContext(ctx context.Context, logger Logger, allocator Allocator) Context {
	return &contextImpl{
		base:      ctx,
		logger:    logger,
		allocator: allocator,
	}
}

func (c *contextImpl) Deadline() (deadline time.Time, ok bool) {
	return c.base.Deadline()
}

func (c *contextImpl) Done() <-chan struct{} {
	return c.base.Done()
}

func (c *contextImpl) Err() error {
	return c.base.Err()
}

func (c *contextImpl) Value(key any) any {
	return c.base.Value(key)
}

func (c *contextImpl) Logger() Logger {
	return c.logger
}

func (c *contextImpl) Alloc() Allocator {
	return c.allocator
}
