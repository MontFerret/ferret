package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// IS_BOOL checks whether value is a boolean value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is boolean, otherwise false.
func IsBool(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return isTypeof(args[0], types.Boolean), nil
}
