package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_BOOL checks whether value is a boolean value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is boolean, otherwise false.
func IsBool(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	return isTypeof(args[0], runtime.TypeBoolean), nil
}
