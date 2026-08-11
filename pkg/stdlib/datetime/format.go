package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE_FORMAT format date according to the given format string.
// @param date {DateTime} Source DateTime object.
// @param format {String} String format.
// @return {String} Formatted date.
func DateFormat(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	dt, format, err := runtime.CastArgs2[runtime.DateTime, runtime.String](arg1, arg2)

	if err != nil {
		return runtime.None, err
	}

	return runtime.String(dt.Format(format.String())), nil
}
