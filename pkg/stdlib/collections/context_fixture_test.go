package collections_test

import "context"

// Host fixtures use the original context so this observes only wrapper polling.
type doneCountingContext struct {
	context.Context
	doneCalls int
}

func (c *doneCountingContext) Done() <-chan struct{} {
	c.doneCalls++

	return c.Context.Done()
}
