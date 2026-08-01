package datetime

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// waitDurationDiff is registered under an unlexable name for compiler-generated
// WAITFOR polling; user code cannot invoke it directly.
func waitDurationDiff(_ context.Context, startValue, endValue runtime.Value) (runtime.Value, error) {
	start, err := runtime.CastDateTime(startValue)
	if err != nil {
		return runtime.None, err
	}

	end, err := runtime.CastDateTime(endValue)
	if err != nil {
		return runtime.None, err
	}

	if start.After(end.Time) {
		start, end = end, start
	}

	return runtime.NewDuration(end.Sub(start.Time)), nil
}
