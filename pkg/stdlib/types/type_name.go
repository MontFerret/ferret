package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// TYPENAME returns the data type name of value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns string representation of a type.
func TypeName(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	return values.NewString(core.Reflect(args[0]).Name()), nil
}
