package utils

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// WAIT pauses the execution for a given period.
// @param timeout {Any} A non-negative value coercible to Duration.
// @return {None} None.
func Wait(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	arg, err := runtime.ToDuration(ctx, arg1)
	if err != nil {
		return runtime.None, err
	}

	if arg < 0 {
		return runtime.None, runtime.Error(runtime.ErrInvalidArgument, "wait duration must not be negative")
	}

	timer := time.NewTimer(time.Duration(arg))
	select {
	case <-ctx.Done():
		timer.Stop()
		return runtime.None, ctx.Err()
	case <-timer.C:
	}

	return runtime.None, nil
}
