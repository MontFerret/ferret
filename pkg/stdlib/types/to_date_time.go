package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TO_DATETIME takes an input value of any type and converts it into the appropriate date time value.
// @param {Any} value - Input value of arbitrary type.
// @return {DateTime} - Parsed date time.
func ToDateTime(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	return runtime.ParseDateTime(args[0].String())
}
