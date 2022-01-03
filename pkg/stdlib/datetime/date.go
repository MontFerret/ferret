package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DATE parses a formatted string and returns DateTime object it represents.
// @param {String} time - String representation of DateTime.
// @param {String} [layout = "2006-01-02T15:04:05Z07:00"] - String layout.
// @return {DateTime} - New DateTime object derived from timeString.
func Date(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 2); err != nil {
		return values.None, err
	}

	if err := core.ValidateType(args[0], types.String); err != nil {
		return values.None, err
	}

	str := args[0].(values.String)
	layout := values.DefaultTimeLayout

	if len(args) > 1 {
		if err := core.ValidateType(args[1], types.String); err != nil {
			return values.None, err
		}

		layout = values.ToString(args[1]).String()
	}

	t, err := time.Parse(layout, str.String())

	if err != nil {
		return values.None, err
	}

	return values.NewDateTime(t), nil
}
