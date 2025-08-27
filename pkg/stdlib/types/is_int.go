package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_INT checks whether value is a int value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is int, otherwise false.
func IsInt(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	return isTypeof(args[0], runtime.TypeInt), nil
}
