package utils

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// WAIT pauses the execution for a given period.
// @param {Int | Float} timeout - Number value which indicates for how long to stop an execution.
func Wait(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.None, nil
	}

	arg := internal.ToInt(args[0])

	timer := time.NewTimer(time.Millisecond * time.Duration(arg))
	select {
	case <-ctx.Done():
		timer.Stop()
		return core.None, ctx.Err()
	case <-timer.C:
	}

	return core.None, nil
}
