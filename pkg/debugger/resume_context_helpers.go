package debugger

import (
	"context"
	"errors"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func earliestDeadline(contexts ...context.Context) (deadline time.Time, ok bool) {
	for _, ctx := range contexts {
		if ctx == nil {
			continue
		}

		next, exists := ctx.Deadline()

		if exists && (!ok || next.Before(deadline)) {
			deadline = next
			ok = true
		}
	}

	return deadline, ok
}

func propagateCancellation(parent context.Context, cancel context.CancelCauseFunc) func() bool {
	if parent == nil {
		return func() bool { return false }
	}

	return context.AfterFunc(parent, func() {
		if errors.Is(parent.Err(), context.DeadlineExceeded) {
			if _, ok := parent.Deadline(); ok {
				return
			}
		}

		cancel(context.Cause(parent))
	})
}

func checkContext(ctx context.Context) error {
	if ctx == nil {
		return runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	return ctx.Err()
}
