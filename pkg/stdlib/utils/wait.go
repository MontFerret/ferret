package utils

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Wait pauses the execution for a given period.
// @param timeout (Float|Int) - Number value which indicates for how long to stop an execution.
func Wait(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, nil
	}

	arg := values.ToInt(args[0])

	time.Sleep(time.Millisecond * time.Duration(arg))

	return values.None, nil
}
