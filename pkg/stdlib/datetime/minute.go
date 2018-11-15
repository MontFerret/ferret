package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DateMinute returns the minute of date as a number.
// @params date (DateTime) - source DateTime.
// @return (Int) - a minute number.
func DateMinute(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.DateTimeType)
	if err != nil {
		return values.None, err
	}

	min := args[0].(values.DateTime).Minute()

	return values.NewInt(int64(min)), nil
}
