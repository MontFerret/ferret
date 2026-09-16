package debugger

import "context"

type (
	// contextObserver captures the initial check before exposing it to a waiter.
	contextObserver struct {
		context.Context
		check func()
	}

	// cancellationCheckpoint injects cancellation at successive validation points.
	cancellationCheckpoint struct {
		context.Context
		cancel    context.CancelFunc
		remaining int
	}
)

func (c *contextObserver) Err() error {
	err := c.Context.Err()
	c.check()

	return err
}

func (c *cancellationCheckpoint) Err() error {
	c.remaining--
	if c.remaining == 0 {
		c.cancel()
	}

	return c.Context.Err()
}
