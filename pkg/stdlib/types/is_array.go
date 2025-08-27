package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_ARRAY checks whether value is an array value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is array, otherwise false.
func IsArray(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertArray(args[0]); err != nil {
		return runtime.False, nil
	}

	return runtime.True, nil
}
