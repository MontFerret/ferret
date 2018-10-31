package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Now returns new DateTime object with Time equal to time.Now().
// @returns (DateTime) - New DateTime object.
func Now(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 0, 0)
	if err != nil {
		return values.None, err
	}

	return values.NewCurrentDateTime(), nil
}
