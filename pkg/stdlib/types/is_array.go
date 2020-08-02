package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// IS_ARRAY checks whether value is an array value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is array, otherwise false.
func IsArray(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return isTypeof(args[0], types.Array), nil
}
