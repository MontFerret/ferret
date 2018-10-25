package utils

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Wait pauses the execution for a given period.
<<<<<<< HEAD
// @param timeout (Float|Int) - Number value which indicates for how long to stop an execution.
func Wait(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
=======
// @param timeout (Int) - Integer value indication for how long to pause.
func Wait(_ context.Context, inputs ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(inputs, 1, 1)
>>>>>>> 9f24172... rewrite comments

	if err != nil {
		return values.None, nil
	}

	arg, err := values.ToInt(args[0])

	if err != nil {
		return values.None, err
	}

	time.Sleep(time.Millisecond * time.Duration(arg))

	return values.None, nil
}
