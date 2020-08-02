package utils

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// WAIT pauses the execution for a given period.
// @param {Int | Float} timeout - Number value which indicates for how long to stop an execution.
func Wait(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, nil
	}

	arg := values.ToInt(args[0])

	timer := time.NewTimer(time.Millisecond * time.Duration(arg))
	select {
	case <-ctx.Done():
		timer.Stop()
		return values.None, ctx.Err()
	case <-timer.C:
	}

	return values.None, nil
}
