package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DateFormat format date according to the given format string.
// @params date (DateTime) - source DateTime object.
// @return (String) - formatted date.
func DateFormat(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.DateTimeType)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], core.StringType)
	if err != nil {
		return values.None, err
	}

	date := args[0].(values.DateTime)
	format := args[1].(values.String).String()

	return values.NewString(date.Format(format)), nil
}
