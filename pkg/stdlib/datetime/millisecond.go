package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DateMillisecond returns the millisecond of date as a number.
// @params date (DateTime) - source DateTime.
// @return (Int) - a millisecond number.
func DateMillisecond(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.DateTime)
	if err != nil {
		return values.None, err
	}

	msec := args[0].(values.DateTime).Nanosecond() / 1000000

	return values.NewInt(msec), nil
}
