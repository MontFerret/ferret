package debugger

import (
	"context"
	"sync"
)

type resumeContext struct {
	context.Context
	run     context.Context
	command context.Context
}

// newResumeContext preserves values installed by before-run hooks while giving
// the current debugger command control over cancellation and value overrides.
func newResumeContext(run, command context.Context) (context.Context, func(), context.CancelCauseFunc) {
	if command == nil || run == command {
		return run, func() {}, nil
	}

	base := context.Background()
	cancelDeadline := func() {}

	if deadline, ok := earliestDeadline(run, command); ok {
		base, cancelDeadline = context.WithDeadline(base, deadline)
	}

	cancelCtx, cancel := context.WithCancelCause(base)
	stopRun := propagateCancellation(run, cancel)
	stopCommand := propagateCancellation(command, cancel)

	ctx := &resumeContext{
		Context: cancelCtx,
		run:     run,
		command: command,
	}

	var once sync.Once

	return ctx, func() {
		once.Do(func() {
			stopRun()
			stopCommand()
			cancel(context.Canceled)
			cancelDeadline()
		})
	}, cancel
}

func (c *resumeContext) Value(key any) any {
	if value := c.Context.Value(key); value != nil {
		return value
	}

	if value := c.command.Value(key); value != nil {
		return value
	}

	if c.run == nil {
		return nil
	}

	return c.run.Value(key)
}
