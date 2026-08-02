package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// IsDuration implements IS_DURATION.
func IsDuration(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeDuration), nil
}
