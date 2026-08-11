package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// IS_DURATION tests whether a value has the Duration type.
// @param value {Any} Input value of arbitrary type.
// @return {Boolean} True when the value is a Duration; otherwise false.
func IsDuration(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeDuration), nil
}
