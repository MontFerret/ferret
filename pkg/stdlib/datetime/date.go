package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DATE parses a formatted string and returns DateTime object it represents.
// @param {String} time - String representation of DateTime.
// @param {String} [layout = "2006-01-02T15:04:05Z07:00"] - String layout.
// @return {DateTime} - New DateTime object derived from timeString.
func Date(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 2); err != nil {
		return core.None, err
	}

	if err := core.AssertString(args[0]); err != nil {
		return core.None, err
	}

	str := args[0].(core.String)
	layout := core.DefaultTimeLayout

	if len(args) > 1 {
		if err := core.AssertString(args[1]); err != nil {
			return core.None, err
		}

		layout = args[1].String()
	}

	t, err := time.Parse(layout, str.String())

	if err != nil {
		return core.None, err
	}

	return core.NewDateTime(t), nil
}
