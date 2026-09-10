package universal

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Portable callbacks run before Native can validate the caller's context.
func checkOptionContext(ctx context.Context) error {
	if ctx == nil {
		return runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	return ctx.Err()
}
