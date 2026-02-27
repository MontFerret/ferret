package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DATE parses a formatted string and returns DateTime object it represents.
// @param {String} time - String representation of DateTime.
// @param {String} [layout = "2006-01-02T15:04:05Z07:00"] - String layout.
// @return {DateTime} - New DateTime object derived from timeString.
func Date(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	str, err := runtime.CastArgAt[runtime.String](args, 0)

	if err != nil {
		return runtime.None, err
	}

	layout := runtime.DefaultTimeLayout

	if len(args) > 1 {
		l, err := runtime.CastArgAt[runtime.String](args, 1)

		if err != nil {
			return runtime.None, err
		}

		layout = string(l)
	}

	t, err := time.Parse(layout, str.String())

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewDateTime(t), nil
}
