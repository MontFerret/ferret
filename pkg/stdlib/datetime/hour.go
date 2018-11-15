package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DateHour returns the hour of date as a number.
// @params date (DateTime) - source DateTime.
// @return (Int) - a hour number.
func DateHour(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.DateTimeType)
	if err != nil {
		return values.None, err
	}

	hour := args[0].(values.DateTime).Hour()

	return values.NewInt(int64(hour)), nil
}
