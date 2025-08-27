package utils

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// WAIT pauses the execution for a given period.
// @param {Int | Float} timeout - Number value which indicates for how long to stop an execution.
func Wait(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, nil
	}

	arg := runtime.ToIntSafe(ctx, args[0])

	timer := time.NewTimer(time.Millisecond * time.Duration(arg))
	select {
	case <-ctx.Done():
		timer.Stop()
		return runtime.None, ctx.Err()
	case <-timer.C:
	}

	return runtime.None, nil
}
