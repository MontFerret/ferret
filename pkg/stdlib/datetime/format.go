package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE_FORMAT format date according to the given format string.
// @param {DateTime} date - Source DateTime object.
// @param {String} format - String format.
// @return {String} - Formatted date.
func DateFormat(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, 2); err != nil {
		return core.None, err
	}

	if err := core.AssertDateTime(args[0]); err != nil {
		return core.None, err
	}

	if err := core.AssertString(args[1]); err != nil {
		return core.None, err
	}

	date := args[0].(core.DateTime)
	format := args[1].(core.String).String()

	return core.NewString(date.Format(format)), nil
}
