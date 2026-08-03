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
func Date(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	if len(args) == 1 {
		return date1(ctx, args[0])
	}

	return date2(ctx, args[0], args[1])
}

func date1(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return date2(ctx, arg1, runtime.NewString(runtime.DefaultTimeLayout))
}

func date2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	str, err := runtime.CastArg[runtime.String](arg1, 0)

	if err != nil {
		return runtime.None, err
	}

	layout, err := runtime.CastArg[runtime.String](arg2, 1)
	if err != nil {
		return runtime.None, err
	}

	t, err := time.Parse(layout.String(), str.String())

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewDateTime(t), nil
}
