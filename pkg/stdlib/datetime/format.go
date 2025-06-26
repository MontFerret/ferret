package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE_FORMAT format date according to the given format string.
// @param {DateTime} date - Source DateTime object.
// @param {String} format - String format.
// @return {String} - Formatted date.
func DateFormat(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 2); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertDateTime(args[0]); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertString(args[1]); err != nil {
		return runtime.None, err
	}

	date := args[0].(runtime.DateTime)
	format := args[1].(runtime.String).String()

	return runtime.NewString(date.Format(format)), nil
}
