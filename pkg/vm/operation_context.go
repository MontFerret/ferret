package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Raw VM cancellation is observed at execution safepoints. Embedding/session
// admission rejects already-canceled work before reaching this layer.
func validateOperationContext(ctx context.Context) error {
	if ctx == nil {
		return runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	return nil
}
