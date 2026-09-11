package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Parse parses a String using an optional Go time layout, defaulting to RFC3339.
// @param time {String} String representation of DateTime.
// @param layout {String} Go time layout.
// @return {DateTime} Parsed DateTime, preserving the parsed offset.
func Parse(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	if len(args) == 1 {
		return parse1(ctx, args[0])
	}

	return parse2(ctx, args[0], args[1])
}

// parse reads an RFC3339 datetime String.
// @param time {String} String representation of DateTime.
// @return {DateTime} Parsed DateTime, preserving the parsed offset.
func parse1(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return parse2(ctx, arg1, runtime.NewString(runtime.DefaultTimeLayout))
}

// parse reads a datetime String using the supplied Go time layout.
// @param time {String} String representation of DateTime.
// @param layout {String} Go time layout.
// @return {DateTime} Parsed DateTime, preserving the parsed offset.
func parse2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
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
