package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// NOW returns new DateTime object with Time equal to time.Now().
// @return {DateTime} - New DateTime object.
func Now(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 0, 0)
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewCurrentDateTime(), nil
}
