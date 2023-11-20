package datetime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DATE_FORMAT format date according to the given format string.
// @param {DateTime} date - Source DateTime object.
// @param {String} format - String format.
// @return {String} - Formatted date.
func DateFormat(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, 2); err != nil {
		return values.None, err
	}

	if err := values.AssertDateTime(args[0]); err != nil {
		return values.None, err
	}

	if err := values.AssertString(args[1]); err != nil {
		return values.None, err
	}

	date := args[0].(values.DateTime)
	format := args[1].(values.String).String()

	return values.NewString(date.Format(format)), nil
}
