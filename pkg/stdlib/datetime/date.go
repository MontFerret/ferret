package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DATE parses a formatted string and returns DateTime object it represents.
// @param {String} time - String representation of DateTime.
// @param {String} [layout = "2006-01-02T15:04:05Z07:00"] - String layout.
// @return {DateTime} - New DateTime object derived from timeString.
func Date(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertString(args[0]); err != nil {
		return runtime.None, err
	}

	str := args[0].(runtime.String)
	layout := runtime.DefaultTimeLayout

	if len(args) > 1 {
		if err := runtime.AssertString(args[1]); err != nil {
			return runtime.None, err
		}

		layout = args[1].String()
	}

	t, err := time.Parse(layout, str.String())

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewDateTime(t), nil
}
