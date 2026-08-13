package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// date parses a formatted string and returns DateTime object it represents.
// @param time {String} String representation of DateTime.
// @param layout {String} String layout.
// @return {DateTime} New DateTime object derived from timeString.
func Date(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	if len(args) == 1 {
		return date1(ctx, args[0])
	}

	return date2(ctx, args[0], args[1])
}

// date parses a formatted string and returns DateTime object it represents.
// @param time {String} String representation of DateTime.
// @return {DateTime} New DateTime object derived from timeString.
func date1(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return date2(ctx, arg1, runtime.NewString(runtime.DefaultTimeLayout))
}

// date parses a formatted string and returns DateTime object it represents.
// @param time {String} String representation of DateTime.
// @param layout {String} String layout.
// @return {DateTime} New DateTime object derived from timeString.
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
