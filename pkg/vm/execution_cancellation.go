package vm

import "context"

type executionCancellation struct {
	ctx  context.Context
	done <-chan struct{}
}

func newExecutionCancellation(ctx context.Context) executionCancellation {
	if ctx == nil {
		return executionCancellation{}
	}

	return executionCancellation{
		ctx:  ctx,
		done: ctx.Done(),
	}
}

// Check observes cancellation only where the VM declares an execution
// safepoint. Code between safepoints remains atomic from the VM's perspective.
func (c executionCancellation) Check() error {
	if c.done == nil {
		return nil
	}

	return c.checkReady()
}

func (c executionCancellation) checkReady() error {
	select {
	case <-c.done:
		return c.ctx.Err()
	default:
		return nil
	}
}
