package math

import "context"

// Cancel during comparator polling without depending on scheduler timing.
type sortCancelContext struct {
	context.Context
	cancel    context.CancelFunc
	remaining int
}

func (ctx *sortCancelContext) Err() error {
	ctx.remaining--
	if ctx.remaining == 0 {
		ctx.cancel()
	}

	return ctx.Context.Err()
}
