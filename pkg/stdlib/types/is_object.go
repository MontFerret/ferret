package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_OBJECT checks whether value is an object value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is object, otherwise false.
func IsObject(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	return isTypeof(args[0], runtime.TypeObject), nil
}
