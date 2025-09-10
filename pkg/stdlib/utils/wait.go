package utils

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// WAIT pauses the execution for a given period.
// @param {Int | Float} timeout - Number value which indicates for how long to stop an execution.
func Wait(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	arg := runtime.ToIntSafe(ctx, arg1)

	timer := time.NewTimer(time.Millisecond * time.Duration(arg))
	select {
	case <-ctx.Done():
		timer.Stop()
		return runtime.None, ctx.Err()
	case <-timer.C:
	}

	return runtime.None, nil
}
