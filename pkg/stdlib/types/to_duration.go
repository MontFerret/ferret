package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ToDuration returns a native duration or parses one from its string representation.
// Numeric values are not interpreted as durations.
// @param {Duration|String} value - A native duration or duration string.
// @return {Duration} - A native duration value.
func ToDuration(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	switch value := arg.(type) {
	case runtime.Duration:
		return value, nil
	case runtime.String:
		return runtime.ParseDuration(value.String())
	default:
		return runtime.None, runtime.TypeErrorOf(arg, runtime.TypeDuration, runtime.TypeString)
	}
}
